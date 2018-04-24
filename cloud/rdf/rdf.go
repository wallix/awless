package rdf

import "fmt"

// Namespaces
const (
	RdfsNS     = "rdfs"
	RdfNS      = "rdf"
	CloudNS    = "cloud"
	CloudRelNS = "cloud-rel"
	CloudOwlNS = "cloud-owl"
	XsdNS      = "xsd"
	NetNS      = "net"
	NetowlNS   = "net-owl"
)

// Existing terms
var (
	RdfsLabel       = fmt.Sprintf("%s:label", RdfsNS)
	RdfsList        = fmt.Sprintf("%s:list", RdfsNS)
	RdfsDefinedBy   = fmt.Sprintf("%s:isDefinedBy", RdfsNS)
	RdfsDataType    = fmt.Sprintf("%s:Datatype", RdfsNS)
	RdfsSeeAlso     = fmt.Sprintf("%s:seeAlso", RdfsNS)
	RdfsLiteral     = fmt.Sprintf("%s:Literal", RdfsNS)
	RdfsClass       = fmt.Sprintf("%s:Class", RdfsNS)
	RdfsSubProperty = fmt.Sprintf("%s:subPropertyOf", RdfsNS)
	RdfsComment     = fmt.Sprintf("%s:comment", RdfsNS)

	RdfType     = fmt.Sprintf("%s:type", RdfNS)
	RdfProperty = fmt.Sprintf("%s:Property", RdfNS)

	XsdString   = fmt.Sprintf("%s:string", XsdNS)
	XsdBoolean  = fmt.Sprintf("%s:boolean", XsdNS)
	XsdInt      = fmt.Sprintf("%s:int", XsdNS)
	XsdDateTime = fmt.Sprintf("%s:dateTime", XsdNS)
)

// Classes
var (
	Grant = fmt.Sprintf("%s:Grant", CloudOwlNS)

	NetFirewallRule    = fmt.Sprintf("%s:FirewallRule", NetowlNS)
	NetRoute           = fmt.Sprintf("%s:Route", NetowlNS)
	CloudGrantee       = fmt.Sprintf("%s:Grantee", CloudOwlNS)
	KeyValue           = fmt.Sprintf("%s:KeyValue", CloudOwlNS)
	DistributionOrigin = fmt.Sprintf("%s:DistributionOrigin", CloudOwlNS)

	Permission = fmt.Sprintf("%s:permission", CloudNS)
	Grantee    = fmt.Sprintf("%s:grantee", CloudNS)

	NetRouteTargets          = fmt.Sprintf("%s:routeTargets", NetNS)
	NetDestinationPrefixList = fmt.Sprintf("%s:routeDestinationPrefixList", NetNS)
)

// Relations
var (
	ParentOf       = fmt.Sprintf("%s:parentOf", CloudRelNS)
	ChildrenOfRel  = "childrenOf"
	ApplyOn        = fmt.Sprintf("%s:applyOn", CloudRelNS)
	DependingOnRel = "dependingOn"
)

var Labels = make(map[string]string)

type rdfProp struct {
	ID, RdfType, RdfsLabel, RdfsDefinedBy, RdfsDataType string
}

type RDFProperties map[string]rdfProp

func (r RDFProperties) Get(prop string) (rdfProp, error) {
	p, ok := r[prop]
	if !ok {
		return rdfProp{}, fmt.Errorf("property '%s' not found", p)
	}
	return p, nil
}

func (r RDFProperties) IsRDFProperty(id string) bool {
	prop, ok := r[id]
	if !ok {
		return false
	}
	return prop.RdfType == RdfProperty
}

func (r RDFProperties) IsRDFSubProperty(id string) bool {
	prop, ok := r[id]
	if !ok {
		return false
	}
	return prop.RdfType == RdfsSubProperty
}

func (r RDFProperties) IsRDFList(prop string) bool {
	p, ok := r[prop]
	if !ok {
		return false
	}
	return p.RdfsDefinedBy == RdfsList
}

func (r RDFProperties) GetRDFId(label string) (string, error) {
	propId, ok := Labels[label]
	if !ok {
		return "", fmt.Errorf("get property id: label '%s' not found", label)
	}
	return propId, nil
}

func (r RDFProperties) GetDataType(prop string) (string, error) {
	p, ok := r[prop]
	if !ok {
		return "", fmt.Errorf("property '%s' not found", p)
	}
	return p.RdfsDataType, nil
}

func (r RDFProperties) GetLabel(prop string) (string, error) {
	p, ok := r[prop]
	if !ok {
		return "", fmt.Errorf("property '%s' not found", p)
	}
	return p.RdfsLabel, nil
}

func (r RDFProperties) GetDefinedBy(prop string) (string, error) {
	p, ok := r[prop]
	if !ok {
		return "", fmt.Errorf("property '%s' not found", p)
	}
	return p.RdfsDefinedBy, nil
}
