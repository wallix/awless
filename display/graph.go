package display

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

// OneResourceOfGraph prints a resource of a RDF graph according to its display properties
func OneResourceOfGraph(w io.Writer, graph *rdf.Graph, resType rdf.ResourceType, resID string, displayer *ResourceDisplayer) error {
	table := NewTable([]*PropertyDisplayer{{Property: "Property", DontTruncate: true}, {Property: "Value", DontTruncate: true}})
	table.MergeIdenticalCells = false
	node, err := node.NewNodeFromStrings(resType.ToRDFType(), resID)
	if err != nil {
		return err
	}

	properties, err := aws.LoadPropertiesFromGraph(graph, node)
	if err != nil {
		return err
	}
	visitedProps := make(map[string]bool)

	for _, propD := range displayer.Properties {
		visitedProps[propD.firstLevelProperty()] = true
		propD.DontTruncate = true
		valueDisplay := propD.display(propD.propertyValue(properties))
		if valueDisplay != "" {
			table.AddValue("Property", propD.displayName())
			table.AddValue("Value", valueDisplay)
		}
	}
	for key, val := range properties {
		if visited, ok := visitedProps[key]; ok && visited {
			continue
		}
		if val != "" {
			table.AddValue("Property", key)
			table.AddValue("Value", fmt.Sprint(val))
		}
		visitedProps[key] = true
	}

	table.SetSortBy("Property", "Value")

	fmt.Fprintf(w, "%s '%s'\n", strings.Title(resType.String()), nameOrID(node, properties))
	table.Fprint(w)
	return nil
}

// SeveralResourcesOfGraph prints a RDF graph with different type of resources according to there display properties
func SeveralResourcesOfGraph(graph *rdf.Graph, displayer *ServiceDisplayer, onlyIDs bool) {
	table := NewTable([]*PropertyDisplayer{{Property: "Type", DontTruncate: true}, {Property: "Name/Id", DontTruncate: true}, {Property: "Property", DontTruncate: true}, {Property: "Value", DontTruncate: true}})
	table.MergeIdenticalCells = true
	for t := range displayer.Resources {
		nodes, err := graph.NodesForType(t)
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
				table.AddValue("Type", t.String())
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
