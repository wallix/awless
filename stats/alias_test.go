package stats

import (
	"reflect"
	"testing"
)

func TestAliasesLoadAndSave(t *testing.T) {
	db, close := newTestDb()
	defer close()

	a, err := db.GetAliases()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(a), 0; got != want {
		t.Fatalf("got %d, want %d\n", got, want)
	}

	db.AddAlias("alias1", "to1")
	db.AddAlias("alias2", "to2")
	db.AddAlias("alias3", "to3")
	db.AddAlias("alias4", "to1")
	db.AddAlias("alias2", "to3")
	db.AddAlias("alias5", "alias2")
	db.AddAlias("alias6", "alias 2") // invalid target
	db.AddAlias("alias 7", "alias2") // invalid alias
	db.AddAlias("alias8", "")        // invalid target
	db.AddAlias("", "alias9")        // invalid alias

	a, err = db.GetAliases()
	if err != nil {
		t.Fatal(err)
	}
	expected := make(Aliases)
	expected["alias1"] = "to1"
	expected["alias3"] = "to3"
	expected["alias4"] = "to1"
	expected["alias2"] = "to3"
	expected["alias5"] = "to3"
	if got, want := a, expected; !reflect.DeepEqual(a, expected) {
		t.Fatalf("got %+v, want %+v\n", got, want)
	}

	db.DeleteAlias([]string{"alias2", "alias3", "alias5"}...)
	a, err = db.GetAliases()
	if err != nil {
		t.Fatal(err)
	}
	expected = make(Aliases)
	expected["alias1"] = "to1"
	expected["alias4"] = "to1"
	if got, want := a, expected; !reflect.DeepEqual(a, expected) {
		t.Fatalf("got %+v, want %+v\n", got, want)
	}
}
