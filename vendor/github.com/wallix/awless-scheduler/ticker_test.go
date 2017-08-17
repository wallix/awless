package main

import (
	"strings"
	"testing"
	"time"

	"github.com/wallix/awless-scheduler/model"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/driver"
)

func TestTicker(t *testing.T) {
	taskStore = createTmpFSStore()
	defer taskStore.Destroy()

	now := time.Now().UTC()

	// never run
	taskStore.Create(&model.Task{
		Content: "#I will never run because I'm to old",
		RunAt:   now.Add(-80 * time.Minute), RevertAt: now,
		Region: "us-west-1",
	})
	taskStore.Create(&model.Task{
		Content: "create instance name=tata",
		RunAt:   now.Add(-5 * time.Minute), RevertAt: now.Add(1 * time.Second),
		Region: "us-west-1",
	})
	taskStore.Create(&model.Task{
		Content: "delete instance id=toto",
		RunAt:   now.Add(-1 * time.Minute),
		Region:  "us-west-1",
	})
	taskStore.Create(&model.Task{
		Content: "create group unexisting=nothing",
		RunAt:   now.Add(-1 * time.Second),
		Region:  "us-west-1",
	})
	taskStore.Create(&model.Task{
		Content: "create subnet cidr=10.0.0.0/24",
		RunAt:   now.Add(2 * time.Second),
		Region:  "us-west-1",
	})
	taskStore.Create(&model.Task{
		Content: "#test will stop before I run",
		RunAt:   now.Add(30 * time.Minute),
		Region:  "us-west-1",
	})

	defaultCompileEnv = newTemplateEnv(func(key string) (template.Definition, bool) {
		if key == "creategroup" {
			return template.Definition{}, false
		}
		return template.Definition{ExtraParams: []string{"id", "name", "cidr"}}, true
	})

	driversFunc = func(region string) (driver.Driver, error) {
		return &happyDriver{}, nil
	}

	tick := newTicker(taskStore, 1*time.Second)
	go tick.start()
	assertEventContainsMsg(t, <-eventc, "success for delete instance id=toto")
	assertEventContainsMsg(t, <-eventc, "success for create instance name=tata")
	assertEventContainsMsg(t, <-eventc, "failure: cannot find template definition for 'creategroup'")
	assertEventContainsMsg(t, <-eventc, "success for delete instance id=tata")
	assertEventContainsMsg(t, <-eventc, "success for create subnet cidr=10.0.0.0/24")
	tick.stop()
	close(eventc)

}

func assertEventContainsMsg(t *testing.T, ev *event, msg string) {
	if !strings.Contains(ev.String(), msg) {
		t.Fatalf("expected '%s' to contain '%s'", ev, msg)
	}
}
