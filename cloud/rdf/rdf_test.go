package rdf

import "testing"

func TestAllLabelsHaveProperties(t *testing.T) {
	for prop, label := range Labels {
		if _, ok := Properties[label]; !ok {
			t.Fatalf("prop '%s' with label '%s' has no corresponding property", prop, label)
		}
	}
}

func TestAllPropertiesHaveLabel(t *testing.T) {
	for label, rdf := range Properties {
		found := false
		for _, v := range Labels {
			if label == v {
				found = true
			}
		}
		if !found && rdf.RdfType != RdfsSubProperty {
			t.Fatalf("rdf prop with label '%s' has no corresponding entry in labels", label)
		}
	}
}
