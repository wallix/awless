package aws

import "testing"

func TestRegionsValid(t *testing.T) {
	if got, want := stringInSlice("eu-west-1", AllRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("us-east-1", AllRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("us-west-1", AllRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("eu-test-1", AllRegions()), false; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	for _, k := range AllRegions() {
		if got, want := IsValidRegion(k), true; got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	}
	if got, want := IsValidRegion("aa-test-10"), false; got != want {
		t.Errorf("got %t, want %t", got, want)
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
