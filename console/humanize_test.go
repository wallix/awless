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
		{stamp: globalNow, expect: "now"},
		{stamp: globalNow.Add(-5 * time.Second), expect: "5 secs"},
		{stamp: globalNow.Add(-1 * time.Minute), expect: "60 secs"},
		{stamp: globalNow.Add(-3 * time.Minute), expect: "3 mins"},
		{stamp: globalNow.Add(-90 * time.Minute), expect: "90 mins"},
		{stamp: globalNow.Add(-3 * time.Hour), expect: "3 hours"},
		{stamp: globalNow.Add(-24 * time.Hour), expect: "24 hours"},
		{stamp: globalNow.Add(-3 * 24 * time.Hour), expect: "3 days"},
		{stamp: globalNow.Add(-3 * 7 * 24 * time.Hour), expect: "3 weeks"},
		{stamp: globalNow.Add(-3 * 30 * 24 * time.Hour), expect: "3 months"},
		{stamp: globalNow.Add(-3 * 365 * 24 * time.Hour), expect: "3 years"},
	}

	for _, tcase := range tcases {
		if got, want := HumanizeTime(tcase.stamp), tcase.expect; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func TestHumanizeStorage(t *testing.T) {
	tcases := []struct {
		from   uint64
		unit   storageUnit
		expect string
	}{
		{from: 3, unit: b, expect: "3B"},
		{from: 300, unit: b, expect: "300B"},
		{from: 3072, unit: b, expect: "3K"},
		{from: 31457280, unit: b, expect: "30M"},
		{from: 31457285, unit: b, expect: "~30M"},
		{from: 2, unit: kb, expect: "2K"},
		{from: 20, unit: kb, expect: "20K"},
		{from: 2048, unit: kb, expect: "2M"},
		{from: 2070, unit: kb, expect: "~2M"},
		{from: 4096, unit: mb, expect: "4G"},
	}

	for _, tcase := range tcases {
		if got, want := HumanizeStorage(tcase.from, tcase.unit), tcase.expect; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}
