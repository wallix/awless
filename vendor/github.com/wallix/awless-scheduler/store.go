package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/wallix/awless-scheduler/model"
)

type store interface {
	Create(tk *model.Task) error
	Remove(id string) error
	GetTasks() ([]*model.Task, error)
	GetFailures() ([]*model.Task, error)
	MarkAsFailed(id string) error
	Cleanup() error
	Destroy() error
}

type fsStore struct {
	mux sync.Mutex

	root, tasksDir, failuresDir string
}

func NewFSStore(root string) (store, error) {
	tasksDir := filepath.Join(root, "tasks")
	failuresDir := filepath.Join(root, "failures")

	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		return nil, fmt.Errorf("cannot make new store: %s", err)
	}

	if err := os.MkdirAll(failuresDir, 0755); err != nil {
		return nil, fmt.Errorf("cannot make new store: %s", err)
	}

	return &fsStore{root: root, tasksDir: tasksDir, failuresDir: failuresDir}, nil
}

func (fs *fsStore) Create(tk *model.Task) error {
	fs.mux.Lock()
	defer fs.mux.Unlock()

	err := ioutil.WriteFile(filepath.Join(fs.tasksDir, tk.AsFilename()), []byte(tk.Content), 0644)
	if err != nil {
		return fmt.Errorf("cannot create task as file: %s", err)
	}
	return nil
}

func (fs *fsStore) GetTasks() ([]*model.Task, error) {
	tasks := make([]*model.Task, 0)

	for _, file := range fs.getTasks() {
		tk, err := New(file)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, tk)
	}

	return tasks, nil
}

func (fs *fsStore) GetFailures() ([]*model.Task, error) {
	tasks := make([]*model.Task, 0)

	for _, file := range fs.getFailures() {
		tk, err := New(file)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, tk)
	}

	return tasks, nil
}

func (fs *fsStore) Remove(id string) error {
	fs.mux.Lock()
	defer fs.mux.Unlock()

	return os.Remove(filepath.Join(fs.tasksDir, id))
}

func (fs *fsStore) MarkAsFailed(id string) error {
	return os.Rename(filepath.Join(fs.tasksDir, id), filepath.Join(fs.failuresDir, id))
}

func (fs *fsStore) Cleanup() error {
	fs.mux.Lock()
	defer fs.mux.Unlock()

	files, _ := filepath.Glob(filepath.Join(fs.root, "*", fmt.Sprintf("*.%s", model.AwlessFileExt)))
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *fsStore) Destroy() error {
	fs.mux.Lock()
	defer fs.mux.Unlock()

	return os.RemoveAll(fs.root)
}

func (fs *fsStore) getTasks() []string {
	fs.mux.Lock()
	defer fs.mux.Unlock()

	return glob(fs.tasksDir)
}

func (fs *fsStore) getFailures() []string {
	fs.mux.Lock()
	defer fs.mux.Unlock()

	return glob(fs.failuresDir)
}

func glob(root string) []string {
	files, err := filepath.Glob(filepath.Join(root, fmt.Sprintf("*.%s", model.AwlessFileExt)))
	if err != nil {
		log.Println(err)
	}
	sort.Strings(files)
	return files
}
