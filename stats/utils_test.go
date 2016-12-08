package stats

import (
	"io/ioutil"
	"os"
)

func newTestDb() (*DB, func()) {
	f, e := ioutil.TempFile(".", "test.db")
	if e != nil {
		panic(e)
	}

	db, err := OpenDB(f.Name())
	if err != nil {
		panic(err)
	}

	todefer := func() {
		os.Remove(f.Name())
		db.Close()
	}

	return db, todefer
}
