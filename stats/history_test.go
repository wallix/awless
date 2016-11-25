package stats

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestSaveCommandHistory(t *testing.T) {
	f, e := ioutil.TempFile(".", "test.db")
	if e != nil {
		t.Fatal(e)
	}
	defer os.Remove(f.Name())

	db, err := OpenDB(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.FlushHistory(); err != nil {
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

	if err := db.FlushHistory(); err != nil {
		t.Fatal(err)
	}

	if lines, err := db.GetHistory(0); err != nil {
		t.Fatal(err)
	} else if got, want := len(lines), 0; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}

}
