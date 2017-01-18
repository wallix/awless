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
