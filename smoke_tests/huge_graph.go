package main

import (
	"fmt"
	"log"
	"net"
	"path/filepath"
	"time"

	"flag"

	"io/ioutil"

	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

var (
	location = flag.String("location", "~/.awless/aws/rdf", "where to generate the graph")
)

func main() {
	flag.Parse()

	_, localhost, _ := net.ParseCIDR("127.0.0.1/32")
	_, subnetcidr, _ := net.ParseCIDR("10.192.24.0/24")
	_, subnet2cidr, _ := net.ParseCIDR("10.20.24.0/24")
	_, subnet2ipv6, _ := net.ParseCIDR("2001:db8::/110")
	routes := []*graph.Route{
		{Destination: subnetcidr, DestinationPrefixListId: "toto", Targets: []*graph.RouteTarget{{Type: graph.InstanceTarget, Ref: "ref_1", Owner: "me"}, {Type: graph.GatewayTarget, Ref: "ref_2"}}},
		{Destination: subnet2cidr, DestinationIPv6: subnet2ipv6, DestinationPrefixListId: "tata", Targets: []*graph.RouteTarget{{Type: graph.NetworkInterfaceTarget, Ref: "ref_3"}}},
	}
	rules := []*graph.FirewallRule{
		{PortRange: graph.PortRange{FromPort: 80, ToPort: 80}, Protocol: "tcp", IPRanges: []*net.IPNet{localhost, subnetcidr}},
		{PortRange: graph.PortRange{FromPort: 1, ToPort: 1024}, Protocol: "udp", IPRanges: []*net.IPNet{subnetcidr}},
	}

	gph := graph.NewGraph()

	for i := 0; i < 2; i++ {
		vpcId := fmt.Sprintf("vpc%d", i)
		gph.AddResource(resourcetest.VPC(vpcId).Prop(properties.ID, vpcId).Build())

		routeId := fmt.Sprintf("%s_route", vpcId)
		gph.AddResource(resourcetest.RouteTable(routeId).Prop(properties.ID, routeId).Prop(properties.Vpc, vpcId).Prop("Routes", routes).Build())

		for j := 0; j < 3; j++ {
			subId := fmt.Sprintf("%ssub%d", vpcId, j)
			gph.AddResource(resourcetest.Subnet(subId).Prop(properties.Vpc, vpcId).Prop(properties.Default, true).Build())
			for k := 0; k < 1000; k++ {
				instId := fmt.Sprintf("%sinst%d", subId, k)

				secGroup1Id := fmt.Sprintf("%s_securitygroup1", instId)
				gph.AddResource(resourcetest.SecurityGroup(secGroup1Id).Prop("InboundRules", rules).Prop("OutboundRules", rules).Prop(properties.Vpc, vpcId).Prop(properties.Launched, time.Now()).Build())
				secGroup2Id := fmt.Sprintf("%s_securitygroup2", instId)
				gph.AddResource(resourcetest.SecurityGroup(secGroup2Id).Prop("InboundRules", rules).Prop("OutboundRules", rules).Prop(properties.Vpc, vpcId).Prop(properties.Launched, time.Now()).Build())

				gph.AddResource(resourcetest.Instance(instId).Prop("Name", instId+"name").Prop(properties.Subnet, subId).Prop(properties.Vpc, vpcId).Prop(properties.Launched, time.Now()).Prop(properties.SecurityGroups, []string{secGroup1Id, secGroup2Id}).Build())
			}
		}
	}

	path := filepath.Join(*location, "infra.triples")
	log.Printf("generating file at %s", path)
	ioutil.WriteFile(path, []byte(gph.MustMarshal()), 0600)
}
