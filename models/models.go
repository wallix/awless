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

func (t *Region) Json() string {
	j, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	return string(j)
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
