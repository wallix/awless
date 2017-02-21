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
	"testing"
	"time"
)

func TestGetSetDatabaseValues(t *testing.T) {
	db, close := newTestDb()
	defer close()

	value, e := db.GetStringValue("mykey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := value, ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	e = db.SetStringValue("mykey", "myvalue")
	if e != nil {
		t.Fatal(e)
	}

	value, e = db.GetStringValue("mykey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := value, "myvalue"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	i, e := db.GetIntValue("myintkey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := i, 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	e = db.SetIntValue("myintkey", 10)
	if e != nil {
		t.Fatal(e)
	}

	i, e = db.GetIntValue("myintkey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := i, 10; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	stamp, e := db.GetTimeValue("mytimekey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := stamp.IsZero(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	now := time.Now()
	e = db.SetTimeValue("mytimekey", now)
	if e != nil {
		t.Fatal(e)
	}

	stamp, e = db.GetTimeValue("mytimekey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := stamp, now; !want.Equal(want) {
		t.Fatalf("got %s, want %s", got, want)
	}
}
