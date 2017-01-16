package database

import (
	"strings"
	"testing"
)

func TestSaveCommandHistory(t *testing.T) {
	db, close := newTestDb()
	defer close()

	if err := db.EmptyHistory(); err != nil {
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

	if err := db.EmptyHistory(); err != nil {
		t.Fatal(err)
	}

	if lines, err := db.GetHistory(0); err != nil {
		t.Fatal(err)
	} else if got, want := len(lines), 0; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}
