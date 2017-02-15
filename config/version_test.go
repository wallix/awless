package config_test

import (
	"testing"

	"github.com/wallix/awless/config"
)

func TestIsUpgradeOrNot(t *testing.T) {
	tcases := []struct {
		current, latest string
		exp             bool
		revert          bool
	}{
		{current: "", latest: "", exp: false, revert: false},
		{current: "1.0", latest: "2.0", exp: false, revert: false},
		{current: "any", latest: "", exp: false, revert: false},
		{current: "1.a.0", latest: "1.b.0", exp: false, revert: false},

		{current: "0.0.0", latest: "0.0.0", exp: false, revert: false},

		{current: "0.0.0", latest: "0.0.1", exp: true, revert: false},
		{current: "0.0.0", latest: "0.1.0", exp: true, revert: false},
		{current: "0.0.0", latest: "0.1.0", exp: true, revert: false},
		{current: "0.0.0", latest: "1.0.0", exp: true, revert: false},

		{current: "0.0.10", latest: "0.0.1", exp: false, revert: true},
		{current: "0.0.10", latest: "0.0.10", exp: false, revert: false},
		{current: "0.12.0", latest: "0.1.0", exp: false, revert: true},
		{current: "0.12.0", latest: "0.12.0", exp: false, revert: false},
		{current: "10.0.0", latest: "9.0.0", exp: false, revert: true},
		{current: "10.0.0", latest: "10.0.0", exp: false, revert: false},

		{current: "0.0.10", latest: "0.0.11", exp: true, revert: false},
		{current: "0.9.0", latest: "0.10.0", exp: true, revert: false},
		{current: "9.0.0", latest: "10.0.0", exp: true, revert: false},

		{current: "0.1.0", latest: "0.0.2", exp: false, revert: true},
		{current: "1.0.0", latest: "0.10.0", exp: false, revert: true},

		{current: "1.1.0", latest: "1.1.1", exp: true, revert: false},
		{current: "2.1.5", latest: "2.2.0", exp: true, revert: false},
	}

	for _, tc := range tcases {
		if got, want := config.IsUpgrade(tc.current, tc.latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %t, want %t", tc.current, tc.latest, got, want)
		}
		if got, want := config.IsUpgrade(tc.latest, tc.current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %t, want %t", tc.latest, tc.current, got, want)
		}

		// with 'v' prefix
		current := "v" + tc.current
		latest := "v" + tc.latest

		if got, want := config.IsUpgrade(current, latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %t, want %t", current, latest, got, want)
		}
		if got, want := config.IsUpgrade(latest, current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %t, want %t", latest, current, got, want)
		}
	}
}
