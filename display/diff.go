package display

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

// FullDiff displays a table of a diff with both resources and properties diffs (inserted and deleted triples)
func FullDiff(diff *rdf.Diff, rootNode *node.Node, cloudService string) {
	table, err := tableFromDiff(diff, rootNode, cloudService)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	table.Fprint(os.Stdout)
}

// ResourceDiff displays a tree view of a diff with only the changed resources
func ResourceDiff(diff *rdf.Diff, rootNode *node.Node) {
	diff.FullGraph().VisitDepthFirst(rootNode, func(g *rdf.Graph, n *node.Node, distance int) {
		var lit *literal.Literal
		diff, err := g.TriplesForSubjectPredicate(n, rdf.DiffPredicate)
		if len(diff) > 0 && err == nil {
			lit, _ = diff[0].Object().Literal()
		}

		var tabs bytes.Buffer
		for i := 0; i < distance; i++ {
			tabs.WriteByte('\t')
		}

		switch lit {
		case rdf.ExtraLiteral:
			color.Set(color.FgGreen)
			fmt.Fprintf(os.Stdout, "+%s%s, %s\n", tabs.String(), rdf.NewResourceTypeFromRdfType(n.Type().String()).String(), n.ID())
			color.Unset()
		case rdf.MissingLiteral:
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stdout, "-%s%s, %s\n", tabs.String(), rdf.NewResourceTypeFromRdfType(n.Type().String()).String(), n.ID())
			color.Unset()
		default:
			fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), rdf.NewResourceTypeFromRdfType(n.Type().String()).String(), n.ID())
		}
	})
}

func tableFromDiff(diff *rdf.Diff, rootNode *node.Node, service string) (*Table, error) {
	table := NewTable([]*PropertyDisplayer{
		{Property: "Type", DontTruncate: true},
		{Property: "Name/Id", DontTruncate: true},
		{Property: "Property", DontTruncate: true},
		{Property: "Value", DontTruncate: true},
	})
	table.MergeIdenticalCells = true

	err := diff.FullGraph().VisitDepthFirstUnique(rootNode, func(g *rdf.Graph, n *node.Node, distance int) error {
		var lit *literal.Literal
		diffTriples, err := g.TriplesForSubjectPredicate(n, rdf.DiffPredicate)
		if len(diffTriples) > 0 && err == nil {
			lit, _ = diffTriples[0].Object().Literal()
		}
		nCommon, nInserted, nDeleted := aws.InitFromRdfNode(n), aws.InitFromRdfNode(n), aws.InitFromRdfNode(n)

		err = nCommon.UnmarshalFromGraph(diff.CommonGraph())
		if err != nil {
			return err
		}

		err = nInserted.UnmarshalFromGraph(diff.InsertedGraph())
		if err != nil {
			return err
		}

		err = nDeleted.UnmarshalFromGraph(diff.DeletedGraph())
		if err != nil {
			return err
		}

		var displayProperties, propsChanges, rNew bool
		var rName string

		switch lit {
		case rdf.ExtraLiteral:
			displayProperties = true
			rNew = true
			rName = nameOrID(nInserted)
		case rdf.MissingLiteral:
			rName = nameOrID(nDeleted)
			table.AddRow(
				rdf.NewResourceTypeFromRdfType(n.Type().String()).String(),
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
			table.AddRow(rdf.NewResourceTypeFromRdfType(n.Type().String()).String(), color.New(color.FgGreen).SprintFunc()("+ "+n.ID().String()))
		}
		return nil
	})
	if err != nil {
		return table, err
	}

	table.SetSortBy("Type", "Name/Id", "Property", "Value")
	return table, nil
}

func addDiffProperties(table *Table, service string, rType rdf.ResourceType, rName string, rNew bool, insertedProps, deletedProps aws.Properties) (bool, error) {
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

func addDiffProperty(table *Table, name, visitedName, value, resourceID string, resourceType rdf.ResourceType, rNew bool, visited map[string]bool) {
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
