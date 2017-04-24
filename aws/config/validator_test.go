package awsconfig

import (
	"testing"
)

func TestRegionsValid(t *testing.T) {
	if got, want := stringInSlice("eu-west-1", allRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("us-east-1", allRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("us-west-1", allRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("eu-test-1", allRegions()), false; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	for _, k := range allRegions() {
		if got, want := IsValidRegion(k), true; got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	}
	if got, want := IsValidRegion("aa-test-10"), false; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

func TestInstanceTypeValid(t *testing.T) {
	tcases := []struct {
		str    string
		expect bool
	}{
		{"t2.micro", true},
		{"m3.large", true},
		{"t.", false},
		{".", false},
		{"a.", false},
	}
	for _, tcase := range tcases {
		if got, want := isValidInstanceType(tcase.str), tcase.expect; got != want {
			t.Errorf("%s: got %t, want %t", tcase.str, got, want)
		}
	}
}

func stringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
