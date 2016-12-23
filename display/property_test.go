package display

import (
	"testing"

	"github.com/wallix/awless/cloud/aws"
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
	properties := aws.Properties{
		"Id":          "propId",
		"Name":        "propName",
		"StringSlice": []interface{}{"str1", "str2", "str3"},
		"ObjectSlice": []interface{}{
			map[string]interface{}{"Key": "objkey1", "Value": "objvalue1"},
			map[string]interface{}{"objkey2": "objvalue2"},
		},
		"Map": map[string]interface{}{"key1": "value1", "key2": "value2"},
	}

	if got, want := propertyValue(properties, "Id"), "propId"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := propertyValue(properties, "Unknown"), ""; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := propertyValue(properties, "StringSlice[]length"), "3"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := propertyValue(properties, "ObjectSlice[]length"), "2"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := propertyValue(properties, "ObjectSlice[].objkey1"), "objvalue1"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := propertyValue(properties, "ObjectSlice[].objkey2"), "objvalue2"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := propertyValue(properties, "ObjectSlice[].nothere"), ""; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := propertyValue(properties, "Map.key1"), "value1"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := propertyValue(properties, "Map.key2"), "value2"; got != want {
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
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	p.ColoredValues = map[string]string{"none": "blue"}
	if got, want := p.display(str), str; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
}
