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
			fmt.Fprintf(os.Stdout, "+%s%s, %s\n", tabs.String(), rdf.ToResourceType(n.Type().String()), n.ID())
			color.Unset()
		case rdf.MissingLiteral:
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stdout, "-%s%s, %s\n", tabs.String(), rdf.ToResourceType(n.Type().String()), n.ID())
			color.Unset()
		default:
			fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), rdf.ToResourceType(n.Type().String()), n.ID())
		}
	})
}

func tableFromDiff(diff *rdf.Diff, rootNode *node.Node, cloudService string) (*Table, error) {
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

		var displayProperties, changedProperties, newResource bool

		switch lit {
		case rdf.ExtraLiteral:
			displayProperties = true
			newResource = true
		case rdf.MissingLiteral:
			table.AddRow(fmt.Sprint(rdf.ToResourceType(n.Type().String())), color.New(color.FgRed).SprintFunc()("- "+n.ID().String()))
		default:
			displayProperties = true
		}
		if displayProperties {
			changedProperties, err = addDiffProperties(table, n, diff, cloudService, newResource)
			if err != nil {
				return err
			}
		}
		if !changedProperties && newResource {
			table.AddRow(fmt.Sprint(rdf.ToResourceType(n.Type().String())), color.New(color.FgGreen).SprintFunc()("+ "+n.ID().String()))
		}
		return nil
	})
	if err != nil {
		return table, err
	}

	table.SetSortBy("Type", "Name/Id", "Property", "Value")
	return table, nil
}

func addDiffProperties(table *Table, n *node.Node, diff *rdf.Diff, cloudService string, newResource bool) (hasChanges bool, err error) {
	insertedG := rdf.NewGraphFromTriples(diff.Inserted())
	insertedProp, err := aws.LoadPropertiesFromGraph(insertedG, n)
	if err != nil {
		return false, err
	}

	deletedG := rdf.NewGraphFromTriples(diff.Deleted())
	deletedProp, err := aws.LoadPropertiesFromGraph(deletedG, n)
	if err != nil {
		return false, err
	}

	visitedInsertedProp, visitedDeletedProp := make(map[string]bool), make(map[string]bool)
	resourceType := rdf.ToResourceType(rdf.ToResourceType(n.Type().String()))

	if serviceD, ok := PropertiesDisplayer.Services[cloudService]; ok && serviceD != nil {
		if resourceD := serviceD.Resources[resourceType]; resourceD != nil {

			for _, prop := range resourceD.Properties {
				if propVal := prop.propertyValue(insertedProp); propVal != "" {
					addDiffProperty(
						table,
						prop.displayName(),
						prop.firstLevelProperty(),
						prop.displayForceColor("+ "+propVal, color.FgGreen),
						n.ID().String(),
						resourceType,
						newResource,
						visitedInsertedProp,
					)
					hasChanges = true
				}
				if propVal := prop.propertyValue(deletedProp); propVal != "" {
					addDiffProperty(
						table,
						prop.displayName(),
						prop.firstLevelProperty(),
						prop.displayForceColor("- "+propVal, color.FgRed),
						n.ID().String(),
						resourceType,
						newResource,
						visitedDeletedProp,
					)
					hasChanges = true
				}
			}
		}
	}

	// Render inserted/deleted properties with no displayer
	for key, val := range insertedProp {
		if visited, ok := visitedInsertedProp[key]; ok && visited {
			continue
		}
		addDiffProperty(
			table,
			key,
			key,
			color.New(color.FgGreen).SprintFunc()("+ "+fmt.Sprint(val)),
			n.ID().String(),
			resourceType,
			newResource,
			visitedInsertedProp,
		)
		hasChanges = true
	}

	for key, val := range deletedProp {
		if visited, ok := visitedDeletedProp[key]; ok && visited {
			continue
		}
		addDiffProperty(
			table,
			key,
			key,
			color.New(color.FgRed).SprintFunc()("- "+fmt.Sprint(val)),
			n.ID().String(),
			resourceType,
			newResource,
			visitedDeletedProp,
		)
		hasChanges = true
	}

	return hasChanges, nil
}

func addDiffProperty(table *Table, name, visitedName, value, resourceID, resourceType string, newResource bool, visited map[string]bool) {
	visited[visitedName] = true
	resourceDisplayF := fmt.Sprint
	if newResource {
		resourceDisplayF = func(i ...interface{}) string { return color.New(color.FgGreen).SprintFunc()("+ " + fmt.Sprint(i...)) }
	}
	table.AddRow(
		resourceType,
		resourceDisplayF(resourceID),
		name,
		value,
	)
}
