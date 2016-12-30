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

func FullDiff(diff *rdf.Diff, rootNode *node.Node, cloudService string) {
	table, err := tableFromDiff(diff, rootNode, cloudService)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	table.Fprint(os.Stdout)
}

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
			fmt.Fprintf(os.Stdout, "+%s%s, %s\n", tabs.String(), n.Type(), n.ID())
			color.Unset()
		case rdf.MissingLiteral:
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stdout, "-%s%s, %s\n", tabs.String(), n.Type(), n.ID())
			color.Unset()
		default:
			fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
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
			table.AddRow(fmt.Sprint(n.Type()), color.New(color.FgRed).SprintFunc()("- "+n.ID().String()))
		default:
			displayProperties = true
		}
		if displayProperties {
			changedProperties, err = addDiffProperties(table, g, n, diff, cloudService, newResource)
			if err != nil {
				return err
			}
		}
		if !changedProperties && newResource {
			table.AddRow(fmt.Sprint(n.Type()), color.New(color.FgGreen).SprintFunc()("+ "+n.ID().String()))
		}
		return nil
	})
	if err != nil {
		return table, err
	}

	table.SetSortBy("Type", "Name/Id", "Property", "Value")
	return table, nil
}

func addDiffProperties(table *Table, g *rdf.Graph, n *node.Node, diff *rdf.Diff, cloudService string, newResource bool) (hasChanges bool, err error) {
	propertiesT, err := g.TriplesForSubjectPredicate(n, rdf.PropertyPredicate)
	if err != nil {
		return false, err
	}
	var resourceD *ResourceDisplayer
	if serviceD, ok := PropertiesDisplayer.Services[cloudService]; ok && serviceD != nil {
		resourceType := rdf.ToResourceType(n.Type().String())
		resourceD = serviceD.Resources[resourceType]
	}

	for _, t := range propertiesT {
		properties := make(aws.Properties)
		prop, err := aws.NewPropertyFromTriple(t)
		if err != nil {
			return hasChanges, err
		}
		properties[prop.Key] = prop.Value

		propD := &PropertyDisplayer{Property: prop.Key}
		if resourceD != nil && resourceD.Properties[prop.Key] != nil {
			propD = resourceD.Properties[prop.Key]
		}
		if diff.HasInsertedTriple(t) {
			hasChanges = true
			resourceDisplayF := fmt.Sprint
			if newResource {
				resourceDisplayF = func(i ...interface{}) string { return color.New(color.FgGreen).SprintFunc()("+ " + fmt.Sprint(i...)) }
			}
			table.AddRow(
				fmt.Sprint(n.Type()),
				resourceDisplayF(n.ID()),
				propD.displayName(),
				propD.displayForceColor("+ "+propertyValue(properties, propD.Property), color.FgGreen),
			)
		}
		if diff.HasDeletedTriple(t) {
			hasChanges = true

			table.AddRow(
				fmt.Sprint(n.Type()),
				fmt.Sprint(n.ID()),
				propD.displayName(),
				propD.displayForceColor("- "+propertyValue(properties, propD.Property), color.FgRed),
			)
		}
	}
	return hasChanges, nil
}
