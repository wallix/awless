package main

import (
	"testing"

	"time"

	"github.com/wallix/awless-scheduler/client"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/driver"
)

func TestTasksAPI(t *testing.T) {
	taskStore = createTmpFSStore()
	defer taskStore.Destroy()

	service, err := NewSchedulerService(routes(), "127.0.0.1:9090", "127.0.0.1:9091", true)
	if err != nil {
		t.Fatal(err)
	}
	defer service.Close()

	go service.Start()

	driversFunc = func(region string) (driver.Driver, error) {
		return &happyDriver{}, nil
	}

	time.Sleep(1 * time.Second)
	schedClient, err := client.New(service.discoveryURL())
	if err != nil {
		t.Fatal(err)
	}

	postTemplate := func(t *testing.T, txt string) {
		if err := schedClient.Post(client.Form{
			Region:   "us-west-1",
			RunIn:    "2m",
			RevertIn: "2h",
			Template: txt,
		}); err != nil {
			t.Fatal(err)
		}
	}

	tplText := "create user name=toto\ncreate user name=tata"

	t.Run("ping service", func(t *testing.T) {
		if err := schedClient.Ping(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("template successfully received", func(t *testing.T) {
		defer taskStore.Cleanup()

		postTemplate(t, tplText)

		tasks, err := schedClient.ListTasks()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := len(tasks), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := tasks[0].Content, tplText; got != want {
			t.Fatalf("got \n%q\nwant\n%q\n", got, want)
		}
	})

	t.Run("listing templates", func(t *testing.T) {
		defer taskStore.Cleanup()

		postTemplate(t, tplText)

		tasks, err := schedClient.ListTasks()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := len(tasks), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := string(tasks[0].Content), tplText; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		if got, want := string(tasks[0].Region), "us-west-1"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("executing task", func(t *testing.T) {
		defer taskStore.Cleanup()

		postTemplate(t, tplText)

		tasks, err := schedClient.ListTasks()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := len(tasks), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		task := tasks[0]

		env := newTemplateEnv(func(key string) (template.Definition, bool) {
			return template.Definition{ExtraParams: []string{"name", "user"}}, true
		})

		if _, err = executeTask(task, &happyDriver{}, env); err != nil {
			t.Fatal(err)
		}

		tasks, err = schedClient.ListTasks()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := len(tasks), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		revertTplText := "delete user name=tata\ndelete user name=toto"
		if got, want := tasks[0].Content, revertTplText; got != want {
			t.Fatalf("got \n%q\nwant\n%q\n", got, want)
		}
	})

	t.Run("fail executing driver", func(t *testing.T) {
		defer taskStore.Cleanup()

		postTemplate(t, tplText)

		tasks, err := schedClient.ListTasks()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := len(tasks), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		task := tasks[0]

		env := newTemplateEnv(func(key string) (template.Definition, bool) {
			return template.Definition{RequiredParams: []string{"name", "user"}}, true
		})
		if _, err := executeTask(task, &failDriver{}, env); err == nil {
			t.Fatal("expected error, got nil")
		}

		tasks, err = schedClient.ListTasks()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := len(tasks), 0; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		fails, err := schedClient.ListFailures()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := len(fails), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		if got, want := fails[0].Content, tplText; got != want {
			t.Fatalf("got \n%q\nwant\n%q\n", got, want)
		}
	})
}
