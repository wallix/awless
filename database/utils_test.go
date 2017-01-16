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
	f, e := ioutil.TempFile(".", "test.db")
	if e != nil {
		panic(e)
	}

	db, err := Open(f.Name())
	if err != nil {
		panic(err)
	}

	todefer := func() {
		os.Remove(f.Name())
		db.Close()
	}

	return db, todefer
}
