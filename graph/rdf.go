package graph

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/badwolf/triple"
	cloudrdf "github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/graph/internal/rdf"
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

func getPropertyValue(gph *rdf.Graph, propObj *triple.Object, prop string) (interface{}, error) {
	rdfProp, ok := cloudrdf.RdfProperties[prop]
	if !ok {
		return "", fmt.Errorf("property '%s' not found", prop)
	}
	definedBy := rdfProp.RdfsDefinedBy
	dataType := rdfProp.RdfsDataType
	switch {
	case definedBy == cloudrdf.RdfsLiteral, (definedBy == cloudrdf.RdfsList) && (dataType == cloudrdf.XsdString):
		val, err := propObj.Literal()
		if err != nil {
			return nil, err
		}
		propVal, err := val.Text()
		if err != nil {
			return nil, err
		}
		switch dataType {
		case cloudrdf.XsdBoolean:
			return strconv.ParseBool(propVal)
		case cloudrdf.XsdInt:
			return strconv.Atoi(propVal)
		case cloudrdf.XsdDateTime:
			t := time.Time{}
			err := t.UnmarshalText([]byte(propVal))
			if err != nil {
				return nil, err
			}
			return t, nil
		default:
			return propVal, nil
		}
	case definedBy == cloudrdf.RdfsList && dataType == cloudrdf.NetFirewallRule:
		node, err := propObj.Node()
		if err != nil {
			return nil, err
		}
		rule := &FirewallRule{}
		err = rule.unmarshalFromTriples(gph, node)
		if err != nil {
			return nil, err
		}
		return rule, nil
	case definedBy == cloudrdf.RdfsClass, dataType == cloudrdf.RdfsClass:
		val, err := propObj.Node()
		if err != nil {
			return nil, err
		}
		return val.ID().String(), nil
	case definedBy == cloudrdf.RdfsList && dataType == cloudrdf.NetRoute:
		node, err := propObj.Node()
		if err != nil {
			return nil, err
		}
		route := &Route{}
		err = route.unmarshalFromTriples(gph, node)
		if err != nil {
			return nil, err
		}
		return route, nil
	case definedBy == cloudrdf.RdfsList && dataType == cloudrdf.Grant:
		node, err := propObj.Node()
		if err != nil {
			return nil, err
		}
		grant := &Grant{}
		err = grant.unmarshalFromTriples(gph, node)
		if err != nil {
			return nil, err
		}
		return grant, nil
	default:
		return "", fmt.Errorf("get property value: %s is neither literal nor class, nor list", definedBy)
	}
}

func extractUniqueLiteralTextFromTriples(triples []*triple.Triple) (string, error) {
	if ln := len(triples); ln != 1 {
		return "", fmt.Errorf("expected unique, got %d", ln)
	}
	lit, err := triples[0].Object().Literal()
	if err != nil {
		return "", err
	}
	return lit.Text()
}

func randomRdfId() string {
	return fmt.Sprintf("%x", rand.Uint32())
}
