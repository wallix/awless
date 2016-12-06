package driver

import "testing"

func TestTokenRegion(t *testing.T) {

	if got, want := UNKNOWN.IsAction(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	if got, want := VERIFY.IsAction(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	if got, want := INSTANCE.IsResource(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := INSTANCE.IsAction(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	if got, want := REF.IsParam(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := REF.IsResource(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
}
