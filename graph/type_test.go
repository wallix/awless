package graph

import (
	"net"
	"testing"
)

func TestFirewallRuleContainsIP(t *testing.T) {
	tcases := []struct {
		nets   []string
		ip     string
		result bool
	}{
		{[]string{}, "89.87.189.250", false},
		{[]string{"89.0.0.0/8"}, "89.87.189.250", true},
		{[]string{"89.0.0.0/16"}, "89.87.189.250", false},
		{[]string{"89.87.0.0/16"}, "89.87.189.250", true},
		{[]string{"89.0.0.0/0"}, "89.87.1", false},
	}

	for i, tcase := range tcases {
		rule := &FirewallRule{}
		for _, n := range tcase.nets {
			_, ipnet, _ := net.ParseCIDR(n)
			rule.IPRanges = append(rule.IPRanges, ipnet)
		}
		if rule.Contains(tcase.ip) != tcase.result {
			t.Fatalf("%d. case %s in %v expected %t", i+1, tcase.ip, tcase.nets, tcase.result)
		}
	}
}

func TestPortRangeContainsPort(t *testing.T) {
	tcases := []struct {
		prange PortRange
		port   int64
		result bool
	}{
		{PortRange{Any: true}, 2373, true},
		{PortRange{Any: false}, 2373, false},
		{PortRange{FromPort: 22}, 22, true},
		{PortRange{ToPort: 22}, 22, true},
		{PortRange{FromPort: 22, ToPort: 22}, 22, true},
		{PortRange{FromPort: 20, ToPort: 22}, 22, true},
		{PortRange{FromPort: 22, ToPort: 25}, 22, true},
		{PortRange{FromPort: 20, ToPort: 25}, 22, true},
		{PortRange{FromPort: 23, ToPort: 25}, 22, false},
		{PortRange{}, 22, false},
	}

	for i, tcase := range tcases {
		if tcase.prange.Contains(tcase.port) != tcase.result {
			t.Fatalf("%d. case %d in %v expected %t", i+1, tcase.port, tcase.prange, tcase.result)
		}
	}
}
