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
	"testing"

	"github.com/google/badwolf/triple/node"
)

func TestResourceTypeToRdfType(t *testing.T) {
	if got, want := Region.ToRDFString(), "/region"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := Region.String(), "region"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	resourceTypes := []ResourceType{Region, Vpc, Subnet, Instance, User, Role, Group, Policy}
	for _, r := range resourceTypes {
		_, err := node.NewType(r.ToRDFString())
		if err != nil {
			t.Fatal(err)
		}
		if got, want := "/"+r.String(), r.ToRDFString(); got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}
