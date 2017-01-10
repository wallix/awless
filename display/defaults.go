package display

import (
	"github.com/fatih/color"
	"github.com/wallix/awless/rdf"
)

var DefaultsColumnDefinitions = map[rdf.ResourceType][]ColumnDefinition{
	rdf.Instance: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "State"},
			ColoredValues:          map[string]color.Attribute{"running": color.FgGreen, "stopped": color.FgRed},
		},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "KeyName", Friendly: "Access Key"},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP"},
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
	},
	rdf.User: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "Arn"},
		StringColumnDefinition{Prop: "Path"},
		StringColumnDefinition{Prop: "PasswordLastUsed"},
	},
	rdf.Role: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "Arn"},
		StringColumnDefinition{Prop: "CreateDate"},
		StringColumnDefinition{Prop: "Path"},
	},
	rdf.Policy: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "Arn"},
		StringColumnDefinition{Prop: "CreateDate"},
		StringColumnDefinition{Prop: "UpdateDate"},
		StringColumnDefinition{Prop: "Path"},
	},
	rdf.Group: []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "Arn"},
		StringColumnDefinition{Prop: "CreateDate"},
		StringColumnDefinition{Prop: "Path"},
	},
}
