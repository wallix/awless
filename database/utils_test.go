/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package database

import (
	"io/ioutil"
	"os"
)

func newTestDb() (*DB, func()) {
	f, e := ioutil.TempDir(".", "test")
	if e != nil {
		panic(e)
	}

	os.Setenv("__AWLESS_HOME", f)

	db, err := current()
	if err != nil {
		panic(err)
	}

	todefer := func() {
		db.Close()
		os.RemoveAll(f)
	}

	return db, todefer
}
