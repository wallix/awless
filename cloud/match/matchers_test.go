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
package match

import (
	"testing"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestMatchers(t *testing.T) {
	tcases := []struct {
		match    cloud.Matcher
		resource cloud.Resource
		expect   bool
	}{
		{match: Property("Inexisting", "empty"), resource: resourcetest.Instance("i1").Build(), expect: false},
		{match: Property("Prop", "value"), resource: resourcetest.Instance("i1").Prop("Prop", "value").Build(), expect: true},
		{match: And(Property("Inexisting", "empty"), Property("Prop", "value")), resource: resourcetest.Instance("i1").Prop("Prop", "value").Build(), expect: false},
		{match: Or(Property("Inexisting", "empty"), Property("Prop", "value")), resource: resourcetest.Instance("i1").Prop("Prop", "value").Build(), expect: true},
		{match: Or(Property("Inexisting1", ""), Property("Inexisting2", "")), resource: resourcetest.Instance("i1").Build(), expect: false},
		{match: And(Property("Prop1", "value1"), Property("Prop2", "value2")), resource: resourcetest.Instance("i1").Prop("Prop1", "value1").Prop("Prop2", "value2").Build(), expect: true},
		{match: And(Property("Prop1", "value1"), Property("Prop2", "value2")), resource: resourcetest.Instance("i1").Prop("Prop1", "value1").Prop("Prop2", "value2").Build(), expect: true},
		{match: Property("Prop", 42).MatchString(), resource: resourcetest.Instance("i1").Prop("Prop", "42").Build(), expect: true},
		{match: Property("Prop", "WithCase").IgnoreCase(), resource: resourcetest.Instance("i1").Prop("Prop", "WITHCASE").Build(), expect: true},
		{match: Property("Prop", "42").IgnoreCase().MatchString(), resource: resourcetest.Instance("i1").Prop("Prop", 42).Build(), expect: true},
		{match: Property("Prop", "inside").Contains(), resource: resourcetest.Instance("i1").Prop("Prop", "Match inside the content").Build(), expect: true},
		{match: Tag("Key", "Val"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: true},
		{match: Tag("Key", "Notthis"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: false},
		{match: TagKey("Key"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: true},
		{match: TagKey("NotThis"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: false},
		{match: TagValue("Val"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: true},
		{match: TagValue("NotThis"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: false},
	}
	for i, tcase := range tcases {
		if got, want := tcase.match.Match(tcase.resource), tcase.expect; got != want {
			t.Fatalf("%d: got %t, want %t", i+1, got, want)
		}
	}
}
