package params_test

import (
	"strings"
	"testing"

	"github.com/wallix/awless/template/params"
)

func TestValidation(t *testing.T) {
	vals := params.Validators{
		"one": params.MaxLengthOf(3),
		"two": params.MinLengthOf(2),
	}

	err := params.Validate(vals, map[string]interface{}{"one": "morethan3", "two": "o"})
	if err == nil {
		t.Fatal("expected error got none")
	}
	msg := err.Error()
	if got, want := msg, "param validation:"; !strings.Contains(got, want) {
		t.Fatalf("expected '%s' to contains: %s", got, want)
	}
	if got, want := msg, "param 'one': expected max length of 3"; !strings.Contains(got, want) {
		t.Fatalf("expected '%s' to contains: %s", got, want)
	}
	if got, want := msg, "param 'two': expected min length of 2"; !strings.Contains(got, want) {
		t.Fatalf("expected '%s' to contains: %s", got, want)
	}
}
