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

package graph

import (
	"reflect"
	"testing"
)

func TestCompareProperties(t *testing.T) {
	props1 := Properties(map[string]interface{}{
		"one":   1,
		"two":   2,
		"three": "3",
		"four":  4,
	})
	props2 := Properties(map[string]interface{}{
		"zero":  0,
		"two":   2,
		"three": "3",
		"four":  "4",
		"five":  "5",
	})

	exp := Properties(map[string]interface{}{"one": 1, "four": 4})
	if got, want := props1.Substract(props2), exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	exp = Properties(map[string]interface{}{"zero": 0, "four": "4", "five": "5"})
	if got, want := props2.Substract(props1), exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
