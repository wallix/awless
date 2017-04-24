package rdf

import "fmt"

type rdfProp struct {
	ID, RdfType, RdfsLabel, RdfsDefinedBy, RdfsDataType string
}

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

	NetFirewallRule = fmt.Sprintf("%s:FirewallRule", NetowlNS)
	NetRoute        = fmt.Sprintf("%s:Route", NetowlNS)
	CloudGrantee    = fmt.Sprintf("%s:Grantee", CloudOwlNS)

	Permission = fmt.Sprintf("%s:permission", CloudNS)
	Grantee    = fmt.Sprintf("%s:grantee", CloudNS)

	NetRouteTargets          = fmt.Sprintf("%s:routeTargets", NetNS)
	NetDestinationPrefixList = fmt.Sprintf("%s:routeDestinationPrefixList", NetNS)
)

// Relations
var (
	ParentOf = fmt.Sprintf("%s:parentOf", CloudRelNS)
	ApplyOn  = fmt.Sprintf("%s:applyOn", CloudRelNS)
)
