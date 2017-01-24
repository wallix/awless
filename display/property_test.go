package display

import (
	"testing"

	"github.com/wallix/awless/graph"
)

func TestPropertyDisplayName(t *testing.T) {
	p := PropertyDisplayer{Property: "prop", Label: "label"}
	if got, want := p.displayName(), "label"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	p = PropertyDisplayer{Property: "prop", Label: ""}
	if got, want := p.displayName(), "prop"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
}

func TestGetPropertyValue(t *testing.T) {
	properties := graph.Properties{
		"Id":          "propId",
		"Name":        "propName",
		"StringSlice": []interface{}{"str1", "str2", "str3"},
		"ObjectSlice": []interface{}{
			map[string]interface{}{"Key": "objkey1", "Value": "objvalue1"},
			map[string]interface{}{"objkey2": "objvalue2"},
		},
		"Map": map[string]interface{}{"key1": "value1", "key2": "value2"},
	}

	if got, want := (&PropertyDisplayer{Property: "Id"}).propertyValue(properties), "propId"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "Id"}).firstLevelProperty(), "Id"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "Unknown"}).propertyValue(properties), ""; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "Unknown"}).firstLevelProperty(), "Unknown"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "StringSlice[]length"}).propertyValue(properties), "3"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "StringSlice[]length"}).firstLevelProperty(), "StringSlice"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "ObjectSlice[]length"}).propertyValue(properties), "2"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "ObjectSlice[]length"}).firstLevelProperty(), "ObjectSlice"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "ObjectSlice[].objkey1"}).propertyValue(properties), "objvalue1"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "ObjectSlice[].objkey1"}).firstLevelProperty(), "ObjectSlice"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "ObjectSlice[].objkey2"}).propertyValue(properties), "objvalue2"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "ObjectSlice[].objkey2"}).firstLevelProperty(), "ObjectSlice"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "ObjectSlice[].nothere"}).propertyValue(properties), ""; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "ObjectSlice[].nothere"}).firstLevelProperty(), "ObjectSlice"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "Map.key1"}).propertyValue(properties), "value1"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "Map.key1"}).firstLevelProperty(), "Map"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "Map.key2"}).propertyValue(properties), "value2"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := (&PropertyDisplayer{Property: "Map.key2"}).firstLevelProperty(), "Map"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
}

func TestPropertyDisplay(t *testing.T) {
	p := PropertyDisplayer{}
	str := "abcdefghijklmnopqrstuvwxyz"
	if got, want := p.display(str), truncateLeft(str, truncateSize); got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	p.TruncateRight = true
	if got, want := p.display(str), truncateRight(str, truncateSize); got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	p.DontTruncate = true
	if got, want := p.display(str), str; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	p.ColoredValues = map[string]string{str: "red"}
	if got, want := p.display(str), str; got != want {
		t.Fatalf("got '%q', want '%q'", got, want)
	}
	p.ColoredValues = map[string]string{"none": "blue"}
	if got, want := p.display(str), str; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
}
