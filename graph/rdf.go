package graph

import (
	"fmt"
	"math/rand"

	cloudrdf "github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

func isRDFProperty(id string) bool {
	rdfProp, ok := cloudrdf.RdfProperties[id]
	if !ok {
		return false
	}
	return rdfProp.RdfType == cloudrdf.RdfProperty
}

func isRDFSubProperty(id string) bool {
	rdfProp, ok := cloudrdf.RdfProperties[id]
	if !ok {
		return false
	}
	return rdfProp.RdfType == cloudrdf.RdfsSubProperty
}

func isRDFList(prop string) bool {
	rdfProp, ok := cloudrdf.RdfProperties[prop]
	if !ok {
		return false
	}
	return rdfProp.RdfsDefinedBy == cloudrdf.RdfsList
}

func getPropertyRDFId(label string) (string, error) {
	propId, ok := cloudrdf.Labels[label]
	if !ok {
		return "", fmt.Errorf("get property id: label '%s' not found", label)
	}
	return propId, nil
}

func getPropertyDataType(prop string) (string, error) {
	rdfProp, ok := cloudrdf.RdfProperties[prop]
	if !ok {
		return "", fmt.Errorf("property '%s' not found", prop)
	}
	return rdfProp.RdfsDataType, nil
}

func getPropertyLabel(prop string) (string, error) {
	rdfProp, ok := cloudrdf.RdfProperties[prop]
	if !ok {
		return "", fmt.Errorf("property '%s' not found", prop)
	}
	return rdfProp.RdfsLabel, nil
}

func getPropertyDefinedBy(prop string) (string, error) {
	rdfProp, ok := cloudrdf.RdfProperties[prop]
	if !ok {
		return "", fmt.Errorf("property '%s' not found", prop)
	}
	return rdfProp.RdfsDefinedBy, nil
}

func getPropertyValue(gph tstore.RDFGraph, propObj tstore.Object, prop string) (interface{}, error) {
	rdfProp, ok := cloudrdf.RdfProperties[prop]
	if !ok {
		return "", fmt.Errorf("property '%s' not found", prop)
	}
	definedBy := rdfProp.RdfsDefinedBy
	dataType := rdfProp.RdfsDataType
	switch {
	case definedBy == cloudrdf.RdfsLiteral, (definedBy == cloudrdf.RdfsList) && (dataType == cloudrdf.XsdString):
		return tstore.ParseLiteral(propObj)
	case definedBy == cloudrdf.RdfsList && dataType == cloudrdf.NetFirewallRule:
		id, ok := propObj.Resource()
		if !ok {
			return nil, fmt.Errorf("get property '%s': object not resource identifier", prop)
		}
		rule := &FirewallRule{}
		err := rule.unmarshalFromTriples(gph, id)
		if err != nil {
			return nil, err
		}
		return rule, nil
	case definedBy == cloudrdf.RdfsClass, dataType == cloudrdf.RdfsClass:
		id, ok := propObj.Resource()
		if !ok {
			return nil, fmt.Errorf("get property '%s': object not resource identifier", prop)
		}
		return id, nil
	case definedBy == cloudrdf.RdfsList && dataType == cloudrdf.NetRoute:
		id, ok := propObj.Resource()
		if !ok {
			return nil, fmt.Errorf("get property '%s': object not resource identifier", prop)
		}
		route := &Route{}
		err := route.unmarshalFromTriples(gph, id)
		if err != nil {
			return nil, err
		}
		return route, nil
	case definedBy == cloudrdf.RdfsList && dataType == cloudrdf.Grant:
		id, ok := propObj.Resource()
		if !ok {
			return nil, fmt.Errorf("get property '%s': object not resource identifier", prop)
		}
		grant := &Grant{}
		err := grant.unmarshalFromTriples(gph, id)
		if err != nil {
			return nil, err
		}
		return grant, nil
	default:
		return "", fmt.Errorf("get property value: %s is neither literal nor class, nor list", definedBy)
	}
}

func extractUniqueLiteralTextFromTriples(triples []tstore.Triple) (string, error) {
	if ln := len(triples); ln != 1 {
		return "", fmt.Errorf("expected unique, got %d", ln)
	}
	return tstore.ParseString(triples[0].Object())
}

func randomRdfId() string {
	return fmt.Sprintf("%x", rand.Uint32())
}
