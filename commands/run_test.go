package commands

import "testing"

func TestIsCSV(t *testing.T) {
	tcases := []struct {
		input string
		exp   bool
	}{
		{input: "aa", exp: false},
		{input: "[aa]", exp: false},
		{input: "[aa,bb", exp: false},
		{input: "aa,bb]", exp: false},
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
