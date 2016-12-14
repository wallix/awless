package stats

import (
	"testing"
	"time"

	"github.com/wallix/awless/cloud/mockup"
)

func init() {
	mockup.InitMockup()
}

func TestOpenDbGeneratesIdForNewDb(t *testing.T) {
	db, close := newTestDb()

	newId, err := db.GetStringValue(AWLESS_ID_KEY)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(newId), 64; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
	close()

	db, close = newTestDb()
	defer close()

	id, _ := db.GetStringValue(AWLESS_ID_KEY)
	if got, want := id, newId; got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

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
