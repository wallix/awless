package store

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/wallix/awless/models"
)

func BuildRegionTree(region string, awsVpcs []*ec2.Vpc, awsSubnets []*ec2.Subnet, awsInstances []*ec2.Instance) *models.Region {
	var vpcs []*models.Vpc
	var subnets []*models.Subnet
	var instances []*models.Instance

	for _, instance := range awsInstances {
		instances = append(instances, models.Convert(instance).(*models.Instance))
	}

	for _, subnet := range awsSubnets {
		subnets = append(subnets, models.Convert(subnet).(*models.Subnet))
	}

	for _, instance := range instances {
		for _, sub := range subnets {
			if sub.Id == instance.SubnetId {
				sub.Instances = append(sub.Instances, instance)
			}
		}
	}

	for _, vpc := range awsVpcs {
		vpcs = append(vpcs, models.Convert(vpc).(*models.Vpc))
	}

	for _, sub := range subnets {
		for _, vpc := range vpcs {
			if vpc.Id == sub.VpcId {
				vpc.Subnets = append(vpc.Subnets, sub)
			}
		}
	}

	return &models.Region{Id: region, Vpcs: vpcs}
}
