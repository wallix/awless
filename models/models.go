package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Region struct {
	Id   string
	Vpcs []*Vpc
}

type Vpc struct {
	Id      string `aws:"VpcId"`
	Subnets []*Subnet
}

type Subnet struct {
	Id        string `aws:"SubnetId"`
	VpcId     string `aws:"VpcId"`
	Instances []*Instance
}

type Instance struct {
	Id        string `aws:"InstanceId"`
	SubnetId  string `aws:"SubnetId"`
	VpcId     string `aws:"VpcId"`
	PublicIp  string `aws:"PublicIpAddress"`
	PrivateIp string `aws:"PrivateIpAddress"`
}

func (r *Region) AllSubnetsAndInstances() ([]*Subnet, []*Instance) {
	var subs []*Subnet
	var insts []*Instance
	for _, vpc := range r.Vpcs {
		for _, sub := range vpc.Subnets {
			subs = append(subs, sub)
			insts = append(insts, sub.Instances...)
		}
	}
	return subs, insts
}

func (r *Region) AllSubnets() []*Subnet {
	var all []*Subnet
	for _, vpc := range r.Vpcs {
		all = append(all, vpc.Subnets...)
	}
	return all
}

func (r *Region) Json() []byte {
	content, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	return content
}

func (t *Region) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("Region: %s, %d VPC(s)\n", t.Id, len(t.Vpcs)))
	for i, vpc := range t.Vpcs {
		buf.WriteString(fmt.Sprintf("\t%d. VPC %s, %d subnet(s)\n", i+1, vpc.Id, len(vpc.Subnets)))
		for j, sub := range vpc.Subnets {
			buf.WriteString(fmt.Sprintf("\t\t%d. Subnet %s, %d instance(s)\n", j+1, sub.Id, len(sub.Instances)))
			for k, inst := range sub.Instances {
				buf.WriteString(fmt.Sprintf("\t\t\t%d. Instance %s\n", k+1, inst.Id))
			}
		}
	}

	return buf.String()
}
