package main

import (
	"log"
	"time"

	"github.com/wallix/awless-scheduler/model"
)

type ticker struct {
	frequency time.Duration
	store     store
	tick      *time.Ticker
}

func newTicker(store store, dur time.Duration) *ticker {
	t := &ticker{frequency: dur, store: store}
	t.tick = time.NewTicker(t.frequency)
	return t
}

func (t *ticker) start() {
	for {
		select {
		case <-t.tick.C:
			if *debug {
				log.Println("tick")
			}
			executables := t.retrieveExecutableTasks()
			for _, s := range executables {
				d, err := driversFunc(s.Region)
				if err != nil {
					log.Println(err)
					continue
				}

				evt := &event{tk: s}
				evt.tpl, evt.err = executeTask(s, d, defaultCompileEnv)
				eventc <- evt
			}
		}
	}
}

func (t *ticker) stop() {
	t.tick.Stop()
}

func (t *ticker) retrieveExecutableTasks() []*model.Task {
	tasks, err := t.store.GetTasks()
	if err != nil {
		log.Println(err)
	}

	var executables []*model.Task
	for _, tk := range tasks {
		if isExecutable(tk) {
			executables = append(executables, tk)
		}
	}

	return executables
}

func isExecutable(tk *model.Task) bool {
	now := time.Now().UTC()
	limit := now.Add(stillExecutable)
	return tk.RunAt.After(limit) && now.After(tk.RunAt)
}
