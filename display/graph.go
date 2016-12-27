package display

import (
	"fmt"
	"os"

	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

// ResourceOfGraph prints a RDF ResourceOfGraph of one type, according to display properties
func ResourceOfGraph(graph *rdf.Graph, resourceType string, properties []*PropertyDisplayer, sortBy []string, onlyIDs bool, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	table := NewTable(properties)

	nodes, err := graph.NodesForType(resourceType)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, node := range nodes {
		nodeProperties, err := aws.LoadPropertiesFromGraph(graph, node)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, prop := range properties {
			table.AddValue(prop.displayName(), propertyValue(nodeProperties, prop.Property))
		}
	}
	table.SetSortBy(sortBy...)
	if onlyIDs {
		table.FprintColumnValues(os.Stdout, "Id", " ")
	} else {
		table.Fprint(os.Stdout)
	}
}

// SeveralResourcesOfGraph prints a RDF graph with different type of resources according to there display properties
func SeveralResourcesOfGraph(graph *rdf.Graph, properties map[string][]*PropertyDisplayer, onlyIDs bool, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	table := NewTable([]*PropertyDisplayer{{Property: "Type"}, {Property: "Name/Id"}, {Property: "Property"}, {Property: "Value"}})
	table.MergeIdenticalCells = true
	for t := range properties {
		nodes, err := graph.NodesForType("/" + t)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, node := range nodes {
			nodeProperties, err := aws.LoadPropertiesFromGraph(graph, node)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			for _, prop := range properties[t] {
				table.AddValue("Type", t)
				table.AddValue("Name/Id", nameOrID(nodeProperties))
				table.AddValue("Property", prop.displayName())
				table.AddValue("Value", propertyValue(nodeProperties, prop.Property))
			}
		}
	}

	if onlyIDs {
		table.FprintColumnValues(os.Stdout, "Name/Id", " ")
	} else {
		table.Fprint(os.Stdout)
	}
}

func nameOrID(properties aws.Properties) string {
	return fmt.Sprint(properties["Id"])
}
