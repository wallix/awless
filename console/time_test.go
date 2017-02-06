package console

import (
	"testing"
	"time"
)

func TestHumanizeTime(t *testing.T) {
	tcases := []struct {
		stamp  time.Time
		expect string
	}{
		{stamp: time.Now(), expect: "now"},
		{stamp: time.Now().Add(-5 * time.Second), expect: "5 seconds ago"},
		{stamp: time.Now().Add(-1 * time.Minute), expect: "60 seconds ago"},
		{stamp: time.Now().Add(-3 * time.Minute), expect: "3 minutes ago"},
		{stamp: time.Now().Add(-90 * time.Minute), expect: "90 minutes ago"},
		{stamp: time.Now().Add(-3 * time.Hour), expect: "3 hours ago"},
		{stamp: time.Now().Add(-24 * time.Hour), expect: "24 hours ago"},
		{stamp: time.Now().Add(-3 * 24 * time.Hour), expect: "3 days ago"},
		{stamp: time.Now().Add(-3 * 7 * 24 * time.Hour), expect: "3 weeks ago"},
		{stamp: time.Now().Add(-3 * 30 * 24 * time.Hour), expect: "3 months ago"},
		{stamp: time.Now().Add(-3 * 365 * 24 * time.Hour), expect: "3 years ago"},
	}

	for _, tcase := range tcases {
		if got, want := humanizeTime(tcase.stamp), tcase.expect; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}
