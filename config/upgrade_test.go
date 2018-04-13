/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package config

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

func TestUpgradeMessaging(t *testing.T) {
	BuildFor = "zip"
	tserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ua := r.Header.Get("User-Agent"); !strings.HasPrefix(ua, "awless-client-"+Version) {
			t.Fatalf("unexpected user-agent: %s", ua)
		}
		w.Write([]byte(`{"URL":"https://github.com/wallix/awless/releases/latest","Version":"1000.0.0"}`))
	}))
	var buff bytes.Buffer
	if err := notifyIfUpgrade(tserver.URL, &buff); err != nil {
		t.Fatal(err)
	}

	exp := fmt.Sprintf("New version 1000.0.0 available. Checkout the latest features at https://github.com/wallix/awless/blob/master/CHANGELOG.md\nRun `wget -O awless-1000.0.0.tar.gz https://github.com/wallix/awless/releases/download/1000.0.0/awless-%s-%s.tar.gz`\n", runtime.GOOS, runtime.GOARCH)
	if got, want := buff.String(), exp; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestSemverUpgradeOrNot(t *testing.T) {
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

		{current: "0.1.9", latest: "0.1.10", exp: true, revert: false},
	}

	for _, tc := range tcases {
		if got, want := IsSemverUpgrade(tc.current, tc.latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %t, want %t", tc.current, tc.latest, got, want)
		}
		if got, want := IsSemverUpgrade(tc.latest, tc.current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %t, want %t", tc.latest, tc.current, got, want)
		}

		// with both 'v' prefix
		current := "v" + tc.current
		latest := "v" + tc.latest

		if got, want := IsSemverUpgrade(current, latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %t, want %t", current, latest, got, want)
		}
		if got, want := IsSemverUpgrade(latest, current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %t, want %t", latest, current, got, want)
		}

		// with current 'v' prefix
		current = "v" + tc.current
		latest = tc.latest

		if got, want := IsSemverUpgrade(current, latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %t, want %t", current, latest, got, want)
		}
		if got, want := IsSemverUpgrade(latest, current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %t, want %t", latest, current, got, want)
		}

		// with latest 'v' prefix
		current = tc.current
		latest = "v" + tc.latest

		if got, want := IsSemverUpgrade(current, latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %t, want %t", current, latest, got, want)
		}
		if got, want := IsSemverUpgrade(latest, current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %t, want %t", latest, current, got, want)
		}
	}
}

func TestCompareSemver(t *testing.T) {
	tcases := []struct {
		current, latest string
		exp, revert     int
	}{
		{current: "0.0.0", latest: "0.0.0", exp: 0, revert: 0},

		{current: "0.0.0", latest: "0.0.1", exp: -1, revert: 1},
		{current: "0.0.0", latest: "0.1.0", exp: -1, revert: 1},
		{current: "0.0.0", latest: "0.1.0", exp: -1, revert: 1},
		{current: "0.0.0", latest: "1.0.0", exp: -1, revert: 1},

		{current: "0.0.10", latest: "0.0.1", exp: 1, revert: -1},
		{current: "0.0.10", latest: "0.0.10", exp: 0, revert: 0},
		{current: "0.12.0", latest: "0.1.0", exp: 1, revert: -1},
		{current: "0.12.0", latest: "0.12.0", exp: 0, revert: 0},
		{current: "10.0.0", latest: "9.0.0", exp: 1, revert: -1},
		{current: "10.0.0", latest: "10.0.0", exp: 0, revert: 0},

		{current: "0.0.10", latest: "0.0.11", exp: -1, revert: 1},
		{current: "0.9.0", latest: "0.10.0", exp: -1, revert: 1},
		{current: "9.0.0", latest: "10.0.0", exp: -1, revert: 1},

		{current: "0.1.0", latest: "0.0.2", exp: 1, revert: -1},
		{current: "1.0.0", latest: "0.10.0", exp: 1, revert: -1},

		{current: "1.1.0", latest: "1.1.1", exp: -1, revert: 1},
		{current: "2.1.5", latest: "2.2.0", exp: -1, revert: 1},
	}

	MustCompareSemver := func(a, b string) int {
		i, err := CompareSemver(a, b)
		if err != nil {
			t.Fatal(err)
		}
		return i
	}

	for _, tc := range tcases {
		if got, want := MustCompareSemver(tc.current, tc.latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %d, want %d", tc.current, tc.latest, got, want)
		}
		if got, want := MustCompareSemver(tc.latest, tc.current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %d, want %d", tc.latest, tc.current, got, want)
		}

		// with both 'v' prefix
		current := "v" + tc.current
		latest := "v" + tc.latest

		if got, want := MustCompareSemver(current, latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %d, want %d", current, latest, got, want)
		}
		if got, want := MustCompareSemver(latest, current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %d, want %d", latest, current, got, want)
		}

		// with current 'v' prefix
		current = "v" + tc.current
		latest = tc.latest

		if got, want := MustCompareSemver(current, latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %d, want %d", current, latest, got, want)
		}
		if got, want := MustCompareSemver(latest, current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %d, want %d", latest, current, got, want)
		}

		// with latest 'v' prefix
		current = tc.current
		latest = "v" + tc.latest

		if got, want := MustCompareSemver(current, latest), tc.exp; got != want {
			t.Fatalf("%s -> %s, got %d, want %d", current, latest, got, want)
		}
		if got, want := MustCompareSemver(latest, current), tc.revert; got != want {
			t.Fatalf("(revert) %s -> %s, got %d, want %d", latest, current, got, want)
		}
	}
}
