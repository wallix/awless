package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	git "github.com/libgit2/git2go"
)

const (
	dirperm os.FileMode = 0700
	filperm os.FileMode = 0600
)

var (
	// ErrNotExist means that the ressource doesn't exists
	ErrNotExist = errors.New("config file does not exist")

	// ErrExist means that the ressource already exists
	ErrExist = errors.New("config file already exist")
)

// Config of the awless cli
type Config struct {
	root string
	Log  log.Logger
	git  struct {
		repo      *git.Repository
		branch    *git.Branch
		signature *git.Signature
		treeoid   *git.Oid
	}
}

// NewConfig create a new config that points to the root
func NewConfig(root string) *Config {
	c := &Config{
		root: root,
		Log: log.Logger{
			Handler: cli.Default,
			Level:   log.DebugLevel,
		},
	}
	c.initRoot()
	c.gitRepoInit()
	//c.gitBranchInit()
	c.gitSignatureInit()
	return c
}

// initRoot initialize the root path of the project
func (c *Config) initRoot() {
	l := c.Log.WithFields(log.Fields{
		"path": c.root,
	})
	l.Debugf("init root")
	info, err := os.Stat(c.root)
	if os.IsNotExist(err) {
		l.Debugf("mkdir")
		err = os.MkdirAll(c.root, dirperm)
		if err != nil {
			l.WithError(err).Fatalf("mkdir")
		}
	} else if !info.IsDir() {
		l.Fatalf("is not a directory")
	}
}

// gitRepoInit initialize the git repository if needed
func (c *Config) gitRepoInit() {
	l := c.Log.WithFields(log.Fields{
		"path": c.root,
	})
	l.Debugf("git repo init")
	repo, err := git.OpenRepository(c.root)
	if err != nil {
		l.Debugf("cannot open git repository, try to init")
		repo, err = git.InitRepository(c.root, false)
		if err != nil {
			l.WithError(err).Fatalf("cannot init git repository")
		}
		f, err := os.OpenFile(path.Join(c.root, ".git/info/exclude"), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		defer f.Close()
		if err != nil {
			l.WithError(err).Fatalf("cannot open ignore rules")
		}
		_, err = f.WriteString("provider\n")
		if err != nil {
			l.WithError(err).Fatalf("%v cannot write ignore rules", f)
		}
	}
	c.git.repo = repo
}

// gitBranchInit initialize the local git branch if needed
func (c *Config) gitBranchInit() {
	l := c.Log.WithFields(log.Fields{
		"path": c.root,
	})
	l.Debugf("git branch init")
	branch, err := c.git.repo.LookupBranch("master", git.BranchLocal)
	if err != nil {
		head, err := c.git.repo.Head()
		if err != nil {
			l.WithError(err).Fatalf("cannot get git HEAD")
		}
		commit, err := c.git.repo.LookupCommit(head.Target())
		if err != nil {
			l.WithError(err).Fatalf("cannot lookup commit of HEAD")
		}
		branch, err = c.git.repo.CreateBranch("master", commit, false)
		if err != nil {
			l.WithError(err).Fatalf("cannot create branch")
		}
	}
	c.git.branch = branch
}

// gitSignatureInit initialize the git signature
func (c *Config) gitSignatureInit() {
	l := c.Log.WithFields(log.Fields{
		"path": c.root,
	})
	l.Debugf("git signature init")
	s, err := c.git.repo.DefaultSignature()
	if err != nil {
		l.WithError(err).Fatalf("cannot get signature")
	}
	c.git.signature = s
}

// Exists returns true if a ressource exists into the path
func (c *Config) Exists(p string) bool {
	p = path.Join(c.root, p)
	l := c.Log.WithFields(log.Fields{
		"path": p,
	})
	l.Debugf("exists")
	_, err := os.Stat(p)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	l.WithError(err).Fatalf("unexpected error")
	return false
}

// Create an object to the given path, returns
func (c *Config) Create(p string, obj interface{}) error {
	return c.save(p, true, obj)
}

// Save an object to the given path
func (c *Config) Save(p string, obj interface{}) {
	c.save(p, false, obj)
}

func (c *Config) save(p0 string, new bool, obj interface{}) error {
	p := path.Join(c.root, p0)
	l := c.Log.WithFields(log.Fields{
		"path": p,
		"type": reflect.TypeOf(obj),
		"obj":  obj,
	})
	l.Debugf("saving")
	dir := path.Dir(p)
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		l.Debugf("mkdir[%s]", dir)
		err = os.MkdirAll(dir, dirperm)
		if err != nil {
			l.WithError(err).Fatalf("mkdir[%s]", dir)
		}
	} else if !info.IsDir() {
		l.Fatalf("path[%s] is not a directory", dir)
	}
	if new {
		info, err = os.Stat(p)
		if err == nil {
			return ErrExist
		}
	}
	file, err := os.Create(p)
	if err != nil {
		l.WithError(err).Fatalf("cannot open")
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(obj)
	if err != nil {
		l.WithError(err).Fatalf("cannot encode")
	}
	c.gitAdd(p0)
	return nil
}

func (c *Config) gitAdd(p string) {
	l := c.Log.WithFields(log.Fields{
		"path": p,
	})
	l.Debugf("git add")
	ignored, err := c.git.repo.IsPathIgnored(p)
	if err != nil {
		l.WithError(err).Fatalf("cannot get ignore path")
	}
	if ignored {
		l.Debugf("do not git add because file is ignored")
		return
	}
	index, err := c.git.repo.Index()
	if err != nil {
		l.WithError(err).Fatalf("cannot get git index")
	}
	err = index.AddByPath(p)
	if err != nil {
		l.WithError(err).Fatalf("cannot git add")
	}
	oid, err := index.WriteTree()
	if err != nil {
		l.WithError(err).Fatalf("cannot write tree")
	}
	c.git.treeoid = oid
	err = index.Write()
	if err != nil {
		l.WithError(err).Fatalf("cannot write index")
	}
}

// Load an obj from the path, return ErrNotExist if the ressource doesn't exists
func (c *Config) Load(p string, obj interface{}) error {
	p = path.Join(c.root, p)
	l := c.Log.WithFields(log.Fields{
		"path": p,
		"type": reflect.TypeOf(obj),
	})
	l.Debugf("loading")
	info, err := os.Stat(p)
	if os.IsNotExist(err) {
		return ErrNotExist
	}
	if info.IsDir() {
		l.Fatalf("cannot open file because is a directory")
	}
	file, err := os.Open(p)
	if err != nil {
		l.WithError(err).Fatalf("cannot open")
	}
	err = json.NewDecoder(file).Decode(obj)
	if err != nil {
		l.WithError(err).Fatalf("cannot decode")
	}
	return nil
}

// List objects on the given path, returns ErrNotExist if the ressource doesn't exists
func (c *Config) List(p string) []string {
	p = path.Join(c.root, p)
	l := c.Log.WithFields(log.Fields{
		"path": p,
	})
	l.Debug("list")
	infos, err := ioutil.ReadDir(p)
	if err != nil {
		l.WithError(err).Fatalf("cannot list directory")
	}
	names := make([]string, len(infos), len(infos))
	for i := range infos {
		names[i] = infos[i].Name()
	}
	return names
}

// Commit the uncommited changes on the configuration
func (c *Config) Commit(msg string, verbose bool) {
	sl, err := c.git.repo.StatusList(nil)
	if err != nil {
		c.Log.WithError(err).Fatalf("status list")
	}
	count, err := sl.EntryCount()
	if err != nil {
		c.Log.WithError(err).Fatalf("entry count")
	}
	if count == 0 {
		c.Log.Error("nothing to commit")
		return
	}
	if verbose {
		for i := 0; i < count; i++ {
			entry, _ := sl.ByIndex(i)
			c.Log.Infof("%v %v", statusToSymbol(entry.Status), entry.HeadToIndex.NewFile.Path)
		}
	}
	if c.git.treeoid == nil {
		i, err := c.git.repo.Index()
		if err != nil {
			c.Log.WithError(err).Fatalf("cannot get index of repository")
		}
		oid, err := i.WriteTree()
		if err != nil {
			c.Log.WithError(err).Fatalf("cannopt write tree")
		}
		c.git.treeoid = oid
	}
	tree, err := c.git.repo.LookupTree(c.git.treeoid)
	if err != nil {
		c.Log.WithError(err).Fatalf("LookupTree error")
	}
	parents := make([]*git.Commit, 0, 1)
	if ref, err := c.git.repo.Head(); err == nil {
		parent, err := c.git.repo.LookupCommit(ref.Target())
		if err != nil {
			c.Log.WithError(err).Fatalf("cannot get latest commit")
		}
		if err == nil {
			parents = append(parents, parent)
		}
	}
	_, err = c.git.repo.CreateCommit("refs/heads/master", c.git.signature, c.git.signature,
		msg, tree, parents...)
	if err != nil {
		c.Log.WithError(err).Fatalf("cannot create commit")
	}
}

func statusToSymbol(status git.Status) string {
	switch status {
	case git.StatusIndexNew:
		return "+"
	case git.StatusIndexModified:
		return "~"
	case git.StatusIndexDeleted:
		return "-"
	}
	return "?"

}
