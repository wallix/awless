package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync"
	"github.com/wallix/awless/sync/repo"
	tstore "github.com/wallix/triplestore"
)

type server struct {
	port string
	gph  *graph.Graph
}

func New(port string) *server {
	return &server{port: port}
}

func (s *server) Start() error {
	g, err := sync.LoadAllGraphs()
	if err != nil {
		return fmt.Errorf("cannot load local graphs: %s", err)
	}

	s.gph = g

	log.Printf("Starting web ui on port %s\n", s.port)
	return http.ListenAndServe(s.port, s.routes())
}

func (s *server) routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/resources/{id}", s.showResourceHandler)
	r.HandleFunc("/resources", s.listResourcesHandler)
	r.HandleFunc("/rdf", s.rdfHandler)
	r.HandleFunc("/", s.homeHandler)
	return r
}

func (s *server) homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("home").Parse(homeTpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) rdfHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(repo.Dir(), fmt.Sprintf("*%s", ".triples"))
	files, _ := filepath.Glob(path)

	var readers []io.Reader
	for _, f := range files {
		reader, err := os.Open(f)
		if err != nil {
			msg := fmt.Sprintf("loading '%s': %s", f, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		readers = append(readers, reader)
	}

	dec := tstore.NewDatasetDecoder(tstore.NewBinaryDecoder, readers...)
	tris, err := dec.Decode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tstore.NewNTriplesEncoder(w).Encode(tris...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) showResourceHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("show").Parse(showResourceTpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resId := mux.Vars(r)["id"]
	res, err := s.gph.FindResource(resId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resource := newResource(res)

	if err := t.Execute(w, resource); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) listResourcesHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("resources").Funcs(template.FuncMap{
		"EscapeForURL": url.PathEscape,
	}).Parse(resourcesTpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resourcesByTypes := make(map[string][]*Resource)

	for _, typ := range aws.ResourceTypes {
		res, err := s.gph.GetAllResources(typ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, r := range res {
			resourcesByTypes[typ] = append(resourcesByTypes[typ], newResource(r))
		}
	}

	if err := t.Execute(w, resourcesByTypes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Resource struct {
	Id, Type   string
	Properties map[string]interface{}
}

func newResource(r *graph.Resource) *Resource {
	return &Resource{Id: r.Id(), Type: r.Type(), Properties: r.Properties}
}

const homeTpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
	<ul>
	<li><a href="/resources">List resources</a></li>
	<li><a href="/rdf">View RDF</a></li>
	</ul>
	</body>
</html>`

const showResourceTpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
	<h2>{{.Type}}: {{.Id}}</h2>
	<ul>
	{{range $name, $val := .Properties}}
	  <li><b>{{$name}}:</b> {{$val}}</li>
        {{end}}
        </ul>
	</body>
</html>`

const resourcesTpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		{{range $type, $resource := .}}
		<h2>{{$type}}</h2>
		<ul>
		  {{range $resource}}
		         {{ $name := index .Properties "Name" }}
			 <li><b>Id:</b> <a href="/resources/{{ EscapeForURL .Id}}">{{.Id}}</a>{{if $name}} {{if (ne (print $name) "")}}, <b>Name:</b> {{$name}}{{end}}{{end}}</li>
		  {{end}}
		</ul>
		{{end}}
	</body>
</html>`
