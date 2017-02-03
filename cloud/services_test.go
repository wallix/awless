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
	}
	for _, tc := range tcases {
		if got, want := PluralizeResource(tc.in), tc.out; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}
