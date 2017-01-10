package display

import (
	"testing"
	"time"
)

func TestHumanizeTime(t *testing.T) {
	if got, want := humanizeTime(time.Now()), "now"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-5*time.Second)), "5 seconds ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-1*time.Minute)), "60 seconds ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-3*time.Minute)), "3 minutes ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-90*time.Minute)), "90 minutes ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-3*time.Hour)), "3 hours ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-24*time.Hour)), "24 hours ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-3*24*time.Hour)), "3 days ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-14*24*time.Hour)), "2 weeks ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-3*30*24*time.Hour)), "3 months ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanizeTime(time.Now().Add(-2.1*365*24*time.Hour)), "2 years ago"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
