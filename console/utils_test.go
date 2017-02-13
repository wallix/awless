package console

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"unicode"
)

func compareJSON(t *testing.T, actual, expected string) {
	var got interface{}
	if err := json.Unmarshal([]byte(actual), &got); err != nil {
		t.Fatal(err)
	}

	var want interface{}
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected same json:\n%s\n\n%s\n", squash(actual), squash(expected))
	}
}

func squash(s string) string {
	squash := func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}

	return strings.Map(squash, s)
}
