package commands

import "testing"

func TestIsCSV(t *testing.T) {
	tcases := []struct {
		input string
		exp   bool
	}{
		{input: "aa", exp: true},
		{input: "[aa]", exp: true},
		{input: "[aa,bb", exp: false},
		{input: "aa,bb]", exp: false},
		{input: "aa,bb", exp: true},
		{input: "[aa,bb]", exp: true},
		{input: "[a1vZ,0k,123,abcd]", exp: true},
		{input: "", exp: false},
		{input: "[]", exp: false},
		{input: ",", exp: false},
		{input: "", exp: false},
		{input: "[abc,2$'ç]", exp: false},
	}
	for i, tcase := range tcases {
		if got, want := isCSV(tcase.input), tcase.exp; got != want {
			t.Fatalf("%d: %s: got %t, want %t", i+1, tcase.input, got, want)
		}
	}
}

func TestJoinSentence(t *testing.T) {
	tcases := []struct {
		in  []string
		exp string
	}{
		{in: []string{""}, exp: ""},
		{in: []string{"", ""}, exp: " and "},
		{in: []string{"one", "two"}, exp: "one and two"},
		{in: []string{"one", "two", "three"}, exp: "one, two and three"},
	}
	for i, tcase := range tcases {
		if got, want := joinSentence(tcase.in), tcase.exp; got != want {
			t.Errorf("%d. got %q, want %q", i+1, got, want)
		}
	}
}
