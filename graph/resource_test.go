package graph

import (
	"reflect"
	"testing"
)

func TestCompareProperties(t *testing.T) {
	props1 := Properties(map[string]interface{}{
		"one":   1,
		"two":   2,
		"three": "3",
		"four":  4,
	})
	props2 := Properties(map[string]interface{}{
		"zero":  0,
		"two":   2,
		"three": "3",
		"four":  "4",
		"five":  "5",
	})

	exp := Properties(map[string]interface{}{"one": 1, "four": 4})
	if got, want := props1.Substract(props2), exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	exp = Properties(map[string]interface{}{"zero": 0, "four": "4", "five": "5"})
	if got, want := props2.Substract(props1), exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
