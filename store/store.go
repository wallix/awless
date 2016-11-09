package store

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type RegionTree struct {
	id   string
	vpcs []*VpcTree
}

type VpcTree struct {
	id      string
	subnets []*SubnetTree
}

type SubnetTree struct {
	id        string
	vpcId     string
	instances []*Instance
}

type Instance struct {
	id       string
	subnetId string
	info     json.RawMessage
}

func BuildRegionTree(region string, vpcs []*ec2.Vpc, subnets []*ec2.Subnet, awsInstances []*ec2.Instance) *RegionTree {
	var vpcTrees []*VpcTree
	var subnetTrees []*SubnetTree
	var instances []*Instance

	for _, instance := range awsInstances {
		instances = append(instances,
			&Instance{id: aws.StringValue(instance.InstanceId), subnetId: aws.StringValue(instance.SubnetId)},
		)
	}

	for _, subnet := range subnets {
		subnetTrees = append(subnetTrees,
			&SubnetTree{id: aws.StringValue(subnet.SubnetId), vpcId: aws.StringValue(subnet.VpcId)},
		)
	}

	for _, instance := range instances {
		for _, sub := range subnetTrees {
			if sub.id == instance.subnetId {
				sub.instances = append(sub.instances, instance)
			}
		}
	}

	for _, vpc := range vpcs {
		vpcTrees = append(vpcTrees, &VpcTree{id: aws.StringValue(vpc.VpcId)})
	}

	for _, sub := range subnetTrees {
		for _, vpc := range vpcTrees {
			if vpc.id == sub.vpcId {
				vpc.subnets = append(vpc.subnets, sub)
			}
		}
	}

	return &RegionTree{id: region, vpcs: vpcTrees}
}

func (t *RegionTree) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("Region: %s, %d VPC(s)\n", t.id, len(t.vpcs)))
	for i, vpc := range t.vpcs {
		buf.WriteString(fmt.Sprintf("\t%d. VPC %s, %d subnet(s)\n", i+1, vpc.id, len(vpc.subnets)))
		for j, sub := range vpc.subnets {
			buf.WriteString(fmt.Sprintf("\t\t%d. Subnet %s, %d instance(s)\n", j+1, sub.id, len(sub.instances)))
			for k, inst := range sub.instances {
				buf.WriteString(fmt.Sprintf("\t\t\t%d. Instance %s\n", k+1, inst.id))
			}
		}
	}

	return buf.String()
}
