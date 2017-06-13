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

package cloud

import "testing"

func TestResourceTypePluralizeName(t *testing.T) {
	tcases := []struct {
		in, out string
	}{
		{in: "region", out: "regions"},
		{in: "vpc", out: "vpcs"},
		{in: "subnet", out: "subnets"},
		{in: "instance", out: "instances"},
		{in: "user", out: "users"},
		{in: "role", out: "roles"},
		{in: "group", out: "groups"},
		{in: "policy", out: "policies"},
		{in: "internetgateway", out: "internetgateways"},
		{in: "repository", out: "repositories"},
		{in: "registry", out: "registries"},
	}
	for _, tc := range tcases {
		if got, want := PluralizeResource(tc.in), tc.out; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func TestResourceTypeSingularizeName(t *testing.T) {
	tcases := []struct {
		in, out string
	}{
		{out: "instance", in: "instance"},
		{out: "region", in: "regions"},
		{out: "vpc", in: "vpcs"},
		{out: "subnet", in: "subnets"},
		{out: "instance", in: "instances"},
		{out: "user", in: "users"},
		{out: "role", in: "roles"},
		{out: "group", in: "groups"},
		{out: "policy", in: "policies"},
		{out: "internetgateway", in: "internetgateways"},
		{out: "repository", in: "repositories"},
		{out: "registry", in: "registries"},
	}
	for _, tc := range tcases {
		if got, want := SingularizeResource(tc.in), tc.out; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}
