package display

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph"
)

// FullDiff displays a table of a diff with both resources and properties diffs (inserted and deleted triples)
func FullDiff(diff *graph.Diff, rootNode *node.Node, cloudService string) {
	table, err := tableFromDiff(diff, rootNode, cloudService)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	table.Fprint(os.Stdout)
}

// ResourceDiff displays a tree view of a diff with only the changed resources
func ResourceDiff(diff *graph.Diff, rootNode *node.Node) {
	diff.FullGraph().Visit(rootNode, func(g *graph.Graph, n *node.Node, distance int) {
		var lit *literal.Literal
		diff, err := g.TriplesInDiff(n)
		if len(diff) > 0 && err == nil {
			lit, _ = diff[0].Object().Literal()
		}

		var tabs bytes.Buffer
		for i := 0; i < distance; i++ {
			tabs.WriteByte('\t')
		}

		var litString string
		if lit != nil {
			litString, _ = lit.Text()
		}

		switch litString {
		case "extra":
			color.Set(color.FgGreen)
			fmt.Fprintf(os.Stdout, "+%s%s, %s\n", tabs.String(), graph.NewResourceType(n.Type()).String(), n.ID())
			color.Unset()
		case "missing":
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stdout, "-%s%s, %s\n", tabs.String(), graph.NewResourceType(n.Type()).String(), n.ID())
			color.Unset()
		default:
			fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), graph.NewResourceType(n.Type()).String(), n.ID())
		}
	})
}

func tableFromDiff(diff *graph.Diff, rootNode *node.Node, service string) (*Table, error) {
	table := NewTable([]*PropertyDisplayer{
		{Property: "Type", DontTruncate: true},
		{Property: "Name/Id", DontTruncate: true},
		{Property: "Property", DontTruncate: true},
		{Property: "Value", DontTruncate: true},
	})
	table.MergeIdenticalCells = true

	err := diff.FullGraph().VisitUnique(rootNode, func(g *graph.Graph, n *node.Node, distance int) error {
		var lit *literal.Literal
		diffTriples, err := g.TriplesInDiff(n)
		if len(diffTriples) > 0 && err == nil {
			lit, _ = diffTriples[0].Object().Literal()
		}
		nCommon, nInserted, nDeleted := graph.InitFromRdfNode(n), graph.InitFromRdfNode(n), graph.InitFromRdfNode(n)

		err = nCommon.UnmarshalFromGraph(&graph.Graph{diff.CommonGraph()})
		if err != nil {
			return err
		}

		err = nInserted.UnmarshalFromGraph(&graph.Graph{diff.InsertedGraph()})
		if err != nil {
			return err
		}

		err = nDeleted.UnmarshalFromGraph(&graph.Graph{diff.DeletedGraph()})
		if err != nil {
			return err
		}

		var displayProperties, propsChanges, rNew bool
		var rName string

		var litString string
		if lit != nil {
			litString, _ = lit.Text()
		}

		switch litString {
		case "extra":
			displayProperties = true
			rNew = true
			rName = nameOrID(nInserted)
		case "missing":
			rName = nameOrID(nDeleted)
			table.AddRow(
				graph.NewResourceType(n.Type()).String(),
				color.New(color.FgRed).SprintFunc()("- "+rName),
			)
		default:
			rName = nameOrID(nCommon)
			displayProperties = true
		}
		if displayProperties {
			propsChanges, err = addDiffProperties(table, service, nCommon.Type(), rName, rNew, nInserted.Properties(), nDeleted.Properties())
			if err != nil {
				return err
			}
		}
		if !propsChanges && rNew {
			table.AddRow(graph.NewResourceType(n.Type()).String(), color.New(color.FgGreen).SprintFunc()("+ "+n.ID().String()))
		}
		return nil
	})
	if err != nil {
		return table, err
	}

	table.SetSortBy("Type", "Name/Id", "Property", "Value")
	return table, nil
}

func addDiffProperties(table *Table, service string, rType graph.ResourceType, rName string, rNew bool, insertedProps, deletedProps graph.Properties) (bool, error) {
	visitedInsertedProp, visitedDeletedProp := make(map[string]bool), make(map[string]bool)
	changes := false

	if serviceD, ok := PropertiesDisplayer.Services[service]; ok && serviceD != nil {
		if resourceD := serviceD.Resources[rType]; resourceD != nil {

			for _, prop := range resourceD.Properties {
				if propVal := prop.propertyValue(insertedProps); propVal != "" {
					addDiffProperty(
						table,
						prop.displayName(),
						prop.firstLevelProperty(),
						prop.displayForceColor("+ "+propVal, color.FgGreen),
						rName,
						rType,
						rNew,
						visitedInsertedProp,
					)
					changes = true
				}
				if propVal := prop.propertyValue(deletedProps); propVal != "" {
					addDiffProperty(
						table,
						prop.displayName(),
						prop.firstLevelProperty(),
						prop.displayForceColor("- "+propVal, color.FgRed),
						rName,
						rType,
						rNew,
						visitedDeletedProp,
					)
					changes = true
				}
			}
		}
	}

	// Render inserted/deleted properties with no displayer
	for key, val := range insertedProps {
		if visited, ok := visitedInsertedProp[key]; ok && visited {
			continue
		}
		addDiffProperty(
			table,
			key,
			key,
			color.New(color.FgGreen).SprintFunc()("+ "+fmt.Sprint(val)),
			rName,
			rType,
			rNew,
			visitedInsertedProp,
		)
		changes = true
	}

	for key, val := range deletedProps {
		if visited, ok := visitedDeletedProp[key]; ok && visited {
			continue
		}
		addDiffProperty(
			table,
			key,
			key,
			color.New(color.FgRed).SprintFunc()("- "+fmt.Sprint(val)),
			rName,
			rType,
			rNew,
			visitedDeletedProp,
		)
		changes = true
	}

	return changes, nil
}

func addDiffProperty(table *Table, name, visitedName, value, resourceID string, resourceType graph.ResourceType, rNew bool, visited map[string]bool) {
	visited[visitedName] = true
	resourceDisplayF := fmt.Sprint
	if rNew {
		resourceDisplayF = func(i ...interface{}) string { return color.New(color.FgGreen).SprintFunc()("+ " + fmt.Sprint(i...)) }
	}
	table.AddRow(
		resourceType.String(),
		resourceDisplayF(resourceID),
		name,
		value,
	)
}
