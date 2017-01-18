package database

import (
	"io/ioutil"
	"os"

	"github.com/wallix/awless/cloud/mocks"
)

func init() {
	mocks.InitServices()
}

func newTestDb() (*DB, func()) {
	f, e := ioutil.TempDir(".", "test")
	if e != nil {
		panic(e)
	}

	os.Setenv("__AWLESS_HOME", f)

	InitDB(true)
	db, closing := Current()

	todefer := func() {
		closing()
		os.RemoveAll(f)
	}

	return db, todefer
}
