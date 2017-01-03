package display

import (
	"fmt"
	"os"

	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

// ResourceOfGraph prints a RDF ResourceOfGraph of one type, according to display properties
func ResourceOfGraph(graph *rdf.Graph, resourceType string, displayer *ResourceDisplayer, sortBy []string, onlyIDs bool) {
	var columnDisplayer []*PropertyDisplayer
	for _, v := range displayer.Properties {
		columnDisplayer = append(columnDisplayer, v)
	}
	table := NewTable(columnDisplayer)

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
		for _, propD := range displayer.Properties {
			table.AddValue(propD.displayName(), propD.propertyValue(nodeProperties))
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
func SeveralResourcesOfGraph(graph *rdf.Graph, displayer *ServiceDisplayer, onlyIDs bool) {
	table := NewTable([]*PropertyDisplayer{{Property: "Type"}, {Property: "Name/Id"}, {Property: "Property"}, {Property: "Value"}})
	table.MergeIdenticalCells = true
	for t := range displayer.Resources {
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
			for _, propD := range displayer.Resources[t].Properties {
				table.AddValue("Type", t)
				table.AddValue("Name/Id", nameOrID(node, nodeProperties))
				table.AddValue("Property", propD.displayName())
				table.AddValue("Value", propD.display(propD.propertyValue(nodeProperties)))
			}
		}
	}

	table.SetSortBy("Type", "Name/Id", "Property", "Value")

	if onlyIDs {
		table.FprintColumnValues(os.Stdout, "Name/Id", " ")
	} else {
		table.Fprint(os.Stdout)
	}
}
