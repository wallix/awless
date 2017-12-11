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

package inspectors

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/graph"
)

type PortScanner struct {
	inbounds   map[string][]*graph.FirewallRule
	applyingOn map[string][]string
}

func (p *PortScanner) Name() string {
	return "port_scanner"
}

func (p *PortScanner) Inspect(g cloud.GraphAPI) error {
	sgroups, err := g.Find(cloud.NewQuery(cloud.SecurityGroup))
	if err != nil {
		return err
	}

	p.inbounds = make(map[string][]*graph.FirewallRule)
	p.applyingOn = make(map[string][]string)
	for _, sg := range sgroups {
		rules := sg.Properties()["InboundRules"]
		switch rules.(type) {
		case []*graph.FirewallRule:
			p.inbounds[sg.Id()] = rules.([]*graph.FirewallRule)
			res, err := g.ResourceRelations(sg, rdf.ApplyOn, false)
			if err != nil {
				return err
			}
			for _, r := range res {
				p.applyingOn[sg.Id()] = append(p.applyingOn[sg.Id()], r.String())
			}

		}
	}

	return nil
}

var allLocalIPs = net.ParseIP("0.0.0.0")

func (p *PortScanner) Print(w io.Writer) {
	for sg, inbounds := range p.inbounds {
		var targets string
		if len(p.applyingOn[sg]) == 0 {
			targets = "nothing"
		} else {
			targets = strings.Join(p.applyingOn[sg], ", ")
		}
		fmt.Fprintf(w, "Securitygroup %s applying on %s: \n", sg, targets)

		var allPermissive bool

		for _, inbound := range inbounds {
			if portRange, prot := inbound.PortRange, inbound.Protocol; portRange.Any == true && prot == "any" {
				var allIps bool
				for _, n := range inbound.IPRanges {
					if n.IP.Equal(allLocalIPs) {
						allIps = true
					}
				}
				if allIps {
					fmt.Fprintf(w, "\tall ports via any protocol for all IPs\n")
				} else {
					fmt.Fprintf(w, "\tall ports via any protocol for IPs: %s\n", inbound.IPRanges)
				}

				allPermissive = true
			}
		}

		if !allPermissive {
			for _, inbound := range inbounds {
				if portRange, prot := inbound.PortRange, inbound.Protocol; prot != "any" {
					if from, to := portRange.FromPort, portRange.ToPort; from == to {
						fmt.Fprintf(w, "\tport %d via %s\n", from, prot)
					} else {
						fmt.Fprintf(w, "\tports %d-%d via %s\n", from, to, prot)
					}
				}
			}
		}
	}
}
