package graph

import (
	"fmt"
	"math/rand"

	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

func getPropertyValue(gph tstore.RDFGraph, propObj tstore.Object, prop string) (interface{}, error) {
	rdfProp, err := rdf.Properties.Get(prop)
	if err != nil {
		return "", err
	}

	definedBy := rdfProp.RdfsDefinedBy
	dataType := rdfProp.RdfsDataType
	switch {
	case definedBy == rdf.RdfsLiteral, (definedBy == rdf.RdfsList) && (dataType == rdf.XsdString):
		return tstore.ParseLiteral(propObj)
	case definedBy == rdf.RdfsList && dataType == rdf.NetFirewallRule:
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
	case definedBy == rdf.RdfsClass, dataType == rdf.RdfsClass:
		id, ok := propObj.Resource()
		if !ok {
			return nil, fmt.Errorf("get property '%s': '%+v' object not resource identifier", prop, propObj)
		}
		return id, nil
	case definedBy == rdf.RdfsList && dataType == rdf.NetRoute:
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
	case definedBy == rdf.RdfsList && dataType == rdf.Grant:
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
	case definedBy == rdf.RdfsList && dataType == rdf.KeyValue:
		id, ok := propObj.Resource()
		if !ok {
			return nil, fmt.Errorf("get property '%s': object not resource identifier", prop)
		}
		kv := &KeyValue{}
		err := kv.unmarshalFromTriples(gph, id)
		if err != nil {
			return nil, err
		}
		return kv, nil
	case definedBy == rdf.RdfsList && dataType == rdf.DistributionOrigin:
		id, ok := propObj.Resource()
		if !ok {
			return nil, fmt.Errorf("get property '%s': object not resource identifier", prop)
		}
		o := &DistributionOrigin{}
		err := o.unmarshalFromTriples(gph, id)
		if err != nil {
			return nil, err
		}
		return o, nil
	default:
		return "", fmt.Errorf("get property value: %s is neither literal nor class, nor list", definedBy)
	}
}

func extractUniqueLiteralTextFromTriples(triples []tstore.Triple) (string, error) {
	if ln := len(triples); ln != 1 {
		return "", fmt.Errorf("expected unique, got %d: %s", ln, triples)
	}
	return tstore.ParseString(triples[0].Object())
}

func randomRdfId() string {
	return fmt.Sprintf("%x", rand.Uint32())
}
