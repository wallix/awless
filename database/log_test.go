package database

import "testing"

func TestAddLogToDatabase(t *testing.T) {
	db, close := newTestDb()
	defer close()

	if err := db.FlushLogs(); err != nil {
		t.Fatal(err)
	}

	if logs, err := db.GetLogs(); err != nil {
		t.Fatal(err)
	} else if got, want := len(logs), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	log1 := "log msg 1"
	log2 := "log msg 2"
	log3 := "log msg 2"
	log4 := "log msg 3"

	if err := db.AddLog(log1); err != nil {
		t.Fatal(err)
	}
	if err := db.AddLog(log2); err != nil {
		t.Fatal(err)
	}
	if err := db.AddLog(log3); err != nil {
		t.Fatal(err)
	}
	if err := db.AddLog(log4); err != nil {
		t.Fatal(err)
	}

	if logs, err := db.GetLogs(); err != nil {
		t.Fatal(err)
	} else if got, want := len(logs), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	} else if got, want := logs[0].Msg, log1; got != want {
		t.Fatalf("got %s, want %s", got, want)
	} else if got, want := logs[0].Hits, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	} else if got, want := logs[1].Msg, log2; got != want {
		t.Fatalf("got %s, want %s", got, want)
	} else if got, want := logs[1].Hits, 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	} else if got, want := logs[2].Msg, log4; got != want {
		t.Fatalf("got %s, want %s", got, want)
	} else if got, want := logs[2].Hits, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if err := db.FlushLogs(); err != nil {
		t.Fatal(err)
	}

	if logs, err := db.GetLogs(); err != nil {
		t.Fatal(err)
	} else if got, want := len(logs), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

}
