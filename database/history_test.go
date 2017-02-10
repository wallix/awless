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
	"strings"
	"testing"
)

func TestSaveCommandHistory(t *testing.T) {
	db, close := newTestDb()
	defer close()

	if err := db.DeleteHistory(); err != nil {
		t.Fatal(err)
	}

	if lines, err := db.GetHistory(0); err != nil {
		t.Fatal(err)
	} else if got, want := len(lines), 0; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}

	if err := db.AddHistoryCommand([]string{"sync"}); err != nil {
		t.Fatal(err)
	}

	if lines, err := db.GetHistory(0); err != nil {
		t.Fatal(err)
	} else if got, want := len(lines), 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	} else if got, want := strings.Join(lines[0].Command, " "), "sync"; got != want {
		t.Fatalf("got %s; want %s", got, want)
	}

	if err := db.DeleteHistory(); err != nil {
		t.Fatal(err)
	}

	if lines, err := db.GetHistory(0); err != nil {
		t.Fatal(err)
	} else if got, want := len(lines), 0; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}
