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

package repo

import (
	"sort"
	"testing"
	"time"
)

func TestSortRev(t *testing.T) {
	revs := []*Rev{
		{Id: "2", Date: time.Now().Add(2 * time.Hour)},
		{Id: "1", Date: time.Now().Add(1 * time.Hour)},
		{Id: "3", Date: time.Now().Add(3 * time.Hour)},
		{Id: "4", Date: time.Now().Add(4 * time.Hour)},
	}

	sort.Sort(revsByDate(revs))

	if got, want := len(revs), 4; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := revs[0].Id, "1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := revs[1].Id, "2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := revs[2].Id, "3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := revs[3].Id, "4"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
