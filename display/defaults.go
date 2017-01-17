package display

import (
	"github.com/fatih/color"
	"github.com/wallix/awless/rdf"
)

var DefaultsColumnDefinitions = map[rdf.ResourceType][]ColumnDefinition{
	rdf.Instance: []ColumnDefinition{
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
	rdf.Vpc: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "IsDefault", Friendly: "Default"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgGreen},
		},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "CidrBlock"},
	},
	rdf.Subnet: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "MapPublicIpOnLaunch", Friendly: "Public VMs"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgYellow}},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "State"},
			ColoredValues:          map[string]color.Attribute{"available": color.FgGreen}},
		StringColumnDefinition{Prop: "CidrBlock"},
		StringColumnDefinition{Prop: "VpcId"},
	},
	rdf.User: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "Arn"},
		StringColumnDefinition{Prop: "Path"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "PasswordLastUsedDate", Friendly: "PasswordLastUsed"}},
	},
	rdf.Role: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "Arn"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
		StringColumnDefinition{Prop: "Path"},
	},
	rdf.Policy: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "Arn"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "UpdateDate"}},
		StringColumnDefinition{Prop: "Path"},
	},
	rdf.Group: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "Arn"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
		StringColumnDefinition{Prop: "Path"},
	},
}
