package repo

import (
	"sort"
	"testing"
	"time"
)

func TestReduceToLastRevOfEachDay(t *testing.T) {
	revs := []*Rev{
		{Id: "1", Date: mustParse("2017-01-18 15:05")},
		{Id: "2", Date: mustParse("2017-01-18 15:09")},
		{Id: "3", Date: mustParse("2017-01-19 09:05")},
		{Id: "4", Date: mustParse("2017-01-19 08:05")},
		{Id: "5", Date: mustParse("2017-01-17 21:05")},
		{Id: "6", Date: mustParse("2017-01-17 10:05")},
	}

	reduced := reduceToLastRevOfEachDay(revs)

	sort.Sort(revsByDate(reduced))

	if got, want := len(reduced), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := reduced[0].Id, "5"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := reduced[1].Id, "2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := reduced[2].Id, "3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func mustParse(s string) time.Time {
	layout := "2006-01-02 15:04"
	t, err := time.Parse(layout, s)
	if err != nil {
		panic(err)
	}
	return t
}
