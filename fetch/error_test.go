package fetch_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/wallix/awless/fetch"
)

func TestAddError(t *testing.T) {
	err := fetch.WrapError()
	err.Add(nil)
	if got, want := err.Any(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	err.Add(fmt.Errorf("anything"))
	if got, want := err.Any(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
}
func TestErrorWrapping(t *testing.T) {
	tcases := []struct {
		err error
		any bool
	}{
		{fmt.Errorf("an error"), true},
		{nil, false},
		{fetch.WrapError(), false},
		{fetch.WrapError(errors.New("one"), errors.New("two")), true},
	}

	for i, tcase := range tcases {
		err := fetch.WrapError(tcase.err)
		if got, want := err.Any(), tcase.any; got != want {
			t.Fatalf("case %d: got %t, want %t. Error: '%s'", i+1, got, want, err)
		}
	}
}
