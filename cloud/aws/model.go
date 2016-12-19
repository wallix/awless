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
		"Id":               "UserId",
		"Name":             "UserName",
		"Arn":              "Arn",
		"Path":             "Path",
		"PasswordLastUsed": "PasswordLastUsed",
	},
	rdf.ROLE: {
		"Id":         "RoleId",
		"Name":       "RoleName",
		"Arn":        "Arn",
		"CreateDate": "CreateDate",
		"Path":       "Path",
	},
	rdf.GROUP: {
		"Id":         "GroupId",
		"Name":       "GroupName",
		"Arn":        "Arn",
		"CreateDate": "CreateDate",
		"Path":       "Path",
	},
	rdf.POLICY: {
		"Id":           "PolicyId",
		"Name":         "PolicyName",
		"Arn":          "Arn",
		"CreateDate":   "CreateDate",
		"UpdateDate":   "UpdateDate",
		"Description":  "Description",
		"IsAttachable": "IsAttachable",
		"Path":         "Path",
	},
}
