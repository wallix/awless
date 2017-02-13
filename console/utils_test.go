package console

import (
	"strings"
	"testing"
	"unicode"
)

func compareJSONString(t *testing.T, expected, actual string) {
	squash := func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}

	if got, want := strings.Map(squash, expected), strings.Map(squash, actual); got != want {
		t.Fatalf("expected same json:\n%q\n\n%q\n", got, want)
	}
}
