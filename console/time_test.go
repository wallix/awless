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
	now := time.Now().UTC()
	tcases := []struct {
		stamp  time.Time
		expect string
	}{
		{stamp: now, expect: "now"},
		{stamp: now.Add(-5 * time.Second), expect: "5 secs"},
		{stamp: now.Add(-1 * time.Minute), expect: "60 secs"},
		{stamp: now.Add(-3 * time.Minute), expect: "3 mins"},
		{stamp: now.Add(-90 * time.Minute), expect: "90 mins"},
		{stamp: now.Add(-3 * time.Hour), expect: "3 hours"},
		{stamp: now.Add(-24 * time.Hour), expect: "24 hours"},
		{stamp: now.Add(-3 * 24 * time.Hour), expect: "3 days"},
		{stamp: now.Add(-3 * 7 * 24 * time.Hour), expect: "3 weeks"},
		{stamp: now.Add(-3 * 30 * 24 * time.Hour), expect: "3 months"},
		{stamp: now.Add(-3 * 365 * 24 * time.Hour), expect: "3 years"},
	}

	for _, tcase := range tcases {
		if got, want := HumanizeTime(tcase.stamp), tcase.expect; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}
