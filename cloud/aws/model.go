package aws

import "github.com/wallix/awless/rdf"

var awsResourcesProperties = map[string]map[string]string{
	rdf.VPC: {
		"Id":        "VpcId",
		"IsDefault": "IsDefault",
		"State":     "State",
		"CidrBlock": "CidrBlock",
	},
	rdf.SUBNET: {
		"Id":                  "SubnetId",
		"VpcId":               "VpcId",
		"MapPublicIpOnLaunch": "MapPublicIpOnLaunch",
		"State":               "State",
		"CidrBlock":           "CidrBlock",
	},
	rdf.INSTANCE: {
		"Id":        "InstanceId",
		"Tags":      "Tags",
		"Type":      "InstanceType",
		"SubnetId":  "SubnetId",
		"VpcId":     "VpcId",
		"PublicIp":  "PublicIpAddress",
		"PrivateIp": "PrivateIpAddress",
		"ImageId":   "ImageId",
		"State":     "State",
	},
	rdf.USER: {
		"Id": "UserId",
	},
	rdf.ROLE: {
		"Id": "RoleId",
	},
	rdf.GROUP: {
		"Id": "GroupId",
	},
	rdf.POLICY: {
		"Id": "PolicyId",
	},
}
