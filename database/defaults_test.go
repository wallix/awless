package database

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	db, close := newTestDb()
	defer close()
	d, err := db.GetDefaults()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := d, make(defaults); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	db.SetDefault("key-1", "value-1")
	db.SetDefault("key-2", "value-2")
	db.SetDefault("key-1", "value-3")

	expected := defaults{
		"key-1": "value-3",
		"key-2": "value-2",
	}

	d, err = db.GetDefaults()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := d, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	v, _ := db.GetDefault("key-1")
	if got, want := v.(string), "value-3"; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	v, _ = db.GetDefault("key-2")
	if got, want := v.(string), "value-2"; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	db.UnsetDefault("key-2")

	expected = defaults{
		"key-1": "value-3",
	}

	d, err = db.GetDefaults()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := d, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	v, ok := db.GetDefault("key-1")
	if got, want := v.(string), "value-3"; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if got, want := ok, true; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	str, ok := db.GetDefaultString("key-1")
	if got, want := str, "value-3"; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if got, want := ok, true; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	_, ok = db.GetDefault("key-2")
	if got, want := ok, false; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	_, ok = db.GetDefaultString("key-2")
	if got, want := ok, false; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestLoadRegion(t *testing.T) {
	f, e := ioutil.TempDir(".", "test")
	if e != nil {
		panic(e)
	}

	os.Setenv("__AWLESS_HOME", f)

	InitDB(true)
	db, closing := MustGetCurrent()

	db.SetDefault(RegionKey, "my-region")
	closing()

	if got, want := MustGetDefaultRegion(), "my-region"; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	os.RemoveAll(f)
}
