package console

import "testing"

func TestTruncate(t *testing.T) {
	if got, want := truncateLeft("ABCDEFGHIJKLM", 5), "...LM"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := truncateLeft("LM", 5), "LM"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := truncateLeft("ABCDEFGHIJKLM", 2), "LM"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}

	if got, want := truncateRight("ABCDEFGHIJKLM", 5), "AB..."; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := truncateRight("LM", 5), "LM"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
	if got, want := truncateRight("ABCDEFGHIJKLM", 2), "AB"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
}
