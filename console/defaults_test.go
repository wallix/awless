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
	"reflect"
	"testing"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
)

func TestGetColumnsDefinitions(t *testing.T) {
	tcases := []struct {
		chosenProperties []string
		resourceType     string
		expectedHeaders  []ColumnDefinition
	}{
		{
			chosenProperties: []string{"id", "name"},
			resourceType:     cloud.Instance,
			expectedHeaders:  []ColumnDefinition{StringColumnDefinition{Prop: properties.ID}, StringColumnDefinition{Prop: properties.Name}},
		},
		{
			chosenProperties: []string{},
			resourceType:     cloud.Instance,
			expectedHeaders:  DefaultsColumnDefinitions[cloud.Instance],
		},
		{
			chosenProperties: []string{"cidr", "Zone", "id", "CIDR"},
			resourceType:     cloud.Subnet,
			expectedHeaders: []ColumnDefinition{
				StringColumnDefinition{Prop: properties.CIDR},
				StringColumnDefinition{Prop: properties.AvailabilityZone, Friendly: "Zone"},
				StringColumnDefinition{Prop: properties.ID},
				StringColumnDefinition{Prop: properties.CIDR},
			},
		},
		{
			chosenProperties: []string{"id", "vpc"},
			resourceType:     cloud.Instance,
			expectedHeaders:  []ColumnDefinition{StringColumnDefinition{Prop: properties.ID}, StringColumnDefinition{Prop: properties.Vpc}},
		},
	}
	for i, tcase := range tcases {
		b := BuildOptions(
			WithRdfType(tcase.resourceType),
			WithColumns(tcase.chosenProperties),
		)
		if got, want := b.columnDefinitions, tcase.expectedHeaders; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %#v,want %#v", i+1, got, want)
		}
	}
}
