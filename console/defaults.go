package console

import (
	"github.com/fatih/color"
	"github.com/wallix/awless/graph"
)

var DefaultsColumnDefinitions = map[graph.ResourceType][]ColumnDefinition{
	//EC2
	graph.Instance: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "SubnetId"},
		StringColumnDefinition{Prop: "Name"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "State"},
			ColoredValues:          map[string]color.Attribute{"running": color.FgGreen, "stopped": color.FgRed},
		},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "KeyName", Friendly: "Access Key"},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "LaunchTime"}},
	},
	graph.Vpc: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "IsDefault", Friendly: "Default"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgGreen},
		},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "CidrBlock"},
	},
	graph.Subnet: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "CidrBlock"},
		StringColumnDefinition{Prop: "AvailabilityZone", Friendly: "Zone"},
		StringColumnDefinition{Prop: "VpcId"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "MapPublicIpOnLaunch", Friendly: "Public VMs"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgYellow}},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "State"},
			ColoredValues:          map[string]color.Attribute{"available": color.FgGreen}},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "DefaultForAz", Friendly: "ZoneDefault"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgGreen},
		},
	},
	graph.SecurityGroup: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "VpcId"},
		FirewallRulesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "InboundRules", Friendly: "Inbound"}},
		FirewallRulesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "OutboundRules", Friendly: "Outbound"}},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		StringColumnDefinition{Prop: "Description", DisableTruncate: true},
	},
	graph.InternetGateway: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		StringColumnDefinition{Prop: "Vpcs", DisableTruncate: true},
	},
	graph.RouteTable: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		StringColumnDefinition{Prop: "VpcId"},
		StringColumnDefinition{Prop: "Main"},
		RoutesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "Routes"}},
	},
	graph.Keypair: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "KeyFingerprint", DisableTruncate: true},
	},
	graph.Volume: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		StringColumnDefinition{Prop: "VolumeType"},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "Size", Friendly: "Size (Gb)"},
		StringColumnDefinition{Prop: "Encrypted"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateTime"}},
		StringColumnDefinition{Prop: "AvailabilityZone"},
	},
	//IAM
	graph.User: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "PasswordLastUsedDate", Friendly: "PasswordLastUsed"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
	},
	graph.Role: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
	},
	graph.Policy: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "UpdateDate"}},
	},
	graph.Group: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
	},
	// S3
	graph.Bucket: []ColumnDefinition{
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		GrantsColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "Grants"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
	},
	graph.Object: []ColumnDefinition{
		StringColumnDefinition{Prop: "Key", TruncateRight: true},
		StringColumnDefinition{Prop: "BucketName"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "ModifiedDate"}},
		StringColumnDefinition{Prop: "OwnerId", TruncateRight: true},
		StringColumnDefinition{Prop: "Size"},
		StringColumnDefinition{Prop: "Class"},
	},
}
