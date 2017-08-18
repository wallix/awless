package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/wallix/awless-scheduler/model"
	"github.com/wallix/awless/aws/driver"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/driver"
)

var (
	discoveryHostport = flag.String("discovery-hostport", "127.0.0.1:8082", "Listening host:port for the discovery service")
	schedulerHostport = flag.String("scheduler-hostport", "127.0.0.1:8083", "Listening host:port for the scheduler service")
	httpMode          = flag.Bool("http-mode", false, "Scheduler service on HTTP")
	tickerFrequency   = flag.Duration("tick-frequency", 1*time.Minute, "ticker frequency to run executable tasks")
	debug             = flag.Bool("debug", false, "print debug messages")
)

var (
	schedulerDir            = filepath.Join(os.Getenv("HOME"), ".awless-scheduler")
	SOCK_ADDR               = filepath.Join(os.Getenv("HOME"), "awless-scheduler.sock")
	minDurationBeforeRevert = 1 * time.Minute
	stillExecutable         = -1 * time.Hour
	eventc                  = make(chan *event)

	taskStore         store
	defaultCompileEnv = awsdriver.DefaultTemplateEnv()
	driversFunc       = func(region string) (driver.Driver, error) { return awsservices.NewDriver(region, "") }
)

func main() {
	flag.Parse()

	var err error
	taskStore, err = NewFSStore(schedulerDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Scheduler home dir: %s", schedulerDir)

	log.Printf("Starting event collector")
	go collectEvents()
	defer close(eventc)

	t := newTicker(taskStore, *tickerFrequency)
	log.Printf("Starting ticker (frequency = %s)", t.frequency)
	go t.start()
	defer t.stop()

	service, err := NewSchedulerService(
		routes(),
		*schedulerHostport,
		*discoveryHostport,
		*httpMode,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
		log.Printf("Service terminated with %s. Cleaning up.", <-sigc)
		service.Close()
	}()

	log.Println(service.Start())
}

type Service struct {
	*http.Server
	listener          *net.UnixListener
	httpMode          bool
	discoveryHostport string
}

func NewSchedulerService(handler http.Handler, serviceHostport, discoveryHostport string, httpMode bool) (*Service, error) {
	s := &http.Server{
		Addr:    serviceHostport,
		Handler: handler,
	}

	service := &Service{Server: s, httpMode: httpMode, discoveryHostport: discoveryHostport}

	if !service.httpMode {
		addr, err := net.ResolveUnixAddr("unix", SOCK_ADDR)
		if err != nil {
			return nil, err
		}
		l, err := net.ListenUnix("unix", addr)
		if err != nil {
			return nil, err
		}
		s.Addr = addr.String()
		service.listener = l
	}

	return service, nil
}

func (s *Service) Start() error {
	go s.startDiscoveryEnpoint()
	log.Printf("Starting scheduler service on %s", s.addr())
	if s.httpMode {
		return s.ListenAndServe()
	}
	return s.Serve(s.listener)
}

func (s *Service) Close() error {
	log.Print("Closing scheduler service")
	return s.Shutdown(context.Background())
}

func (s *Service) addr() string {
	if s.httpMode {
		u := url.URL{Host: s.Addr}
		u.Scheme = "http"
		return u.String()
	}
	return s.Addr
}

func (s *Service) discoveryURL() string {
	u := url.URL{Host: s.discoveryHostport}
	u.Scheme = "http"
	return u.String()
}

func (s *Service) startDiscoveryEnpoint() {
	started := time.Now()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v := model.ServiceInfo{
			TickerFrequency: (*tickerFrequency).String(),
			Uptime:          time.Since(started).String(),
			ServiceAddr:     s.addr(),
			UnixSockMode:    !s.httpMode,
		}
		b, err := json.MarshalIndent(v, "", " ")
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot marshal json for discovery service: %s", err), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	})

	log.Printf("Starting HTTP discovery service on %s", s.discoveryURL())
	log.Fatal(http.ListenAndServe(s.discoveryHostport, nil))
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("scheduler up!"))
	})
	mux.HandleFunc("/tasks", tasks)
	mux.HandleFunc("/failures", listFailures)

	return mux
}

func tasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createTask(w, r)
		return
	} else if r.Method == http.MethodGet {
		listTasks(w, r)
		return
	}
	http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	return
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := taskStore.GetTasks()
	b, err := marshalTasks(tasks)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func listFailures(w http.ResponseWriter, r *http.Request) {
	tasks, err := taskStore.GetFailures()
	b, err := marshalTasks(tasks)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func marshalTasks(tasks []*model.Task) ([]byte, error) {
	sort.Slice(tasks, func(i int, j int) bool { return !tasks[i].RunAt.Before(tasks[j].RunAt) })

	b, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return nil, err
	}

	return b, nil
}

func createTask(w http.ResponseWriter, r *http.Request) {
	if *debug {
		log.Println(r.URL.String())
	}
	region := r.FormValue("region")
	if region == "" {
		log.Println("missing region")
		http.Error(w, "missing region", http.StatusBadRequest)
		return
	}
	runAt, err := getTimeParam(r.FormValue("run"), time.Now().UTC())
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid duration for 'run' param", http.StatusBadRequest)
		return
	}
	revertAt, err := getTimeParam(r.FormValue("revert"), time.Time{})
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid duration for 'revert' param", http.StatusBadRequest)
		return
	}
	if !revertAt.IsZero() && revertAt.Sub(runAt).Seconds() < minDurationBeforeRevert.Seconds() {
		err = fmt.Errorf("revert time is less that %s before run time", minDurationBeforeRevert)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	tplTxt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	tpl, err := template.Parse(string(tplTxt))
	if err != nil {
		errMsg := fmt.Sprintf("cannot parse template: %s", err)
		log.Println(errMsg)
		log.Printf("body was '%s'", string(tplTxt))
		http.Error(w, errMsg, http.StatusUnprocessableEntity)
		return
	}

	env := awsdriver.DefaultTemplateEnv()
	_, _, err = template.Compile(tpl, env)

	if err != nil {
		errMsg := fmt.Sprintf("cannot compile template: %s", err)
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusUnprocessableEntity)
		return
	}
	d, err := driversFunc(region)
	if err != nil {
		errMsg := fmt.Sprintf("cannot init drivers for dryrun: %s", err)
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	env.Driver = d

	if err = tpl.DryRun(env); err != nil {
		errMsg := fmt.Sprintf("cannot dryrun template: %s", err)
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusUnprocessableEntity)
		return
	}

	tk := &model.Task{Content: string(tplTxt), RunAt: runAt, RevertAt: revertAt, Region: region}

	if err := taskStore.Create(tk); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTimeParam(param string, defaultTime time.Time) (time.Time, error) {
	if param == "" {
		return defaultTime, nil
	}

	dur, err := time.ParseDuration(param)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().UTC().Add(dur), nil
}
