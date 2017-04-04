package triplestore

import (
	"fmt"
	"strings"
)

type XsdType string

var (
	XsdString   = XsdType("xsd:string")
	XsdBoolean  = XsdType("xsd:boolean")
	XsdInteger  = XsdType("xsd:integer")
	XsdDateTime = XsdType("xsd:dateTime")
)

const XMLSchemaNamespace = "http://www.w3.org/2001/XMLSchema"

func (x XsdType) NTriplesNamespaced() string {
	splits := strings.Split(string(x), ":")
	if len(splits) != 2 {
		return string(x)
	}

	return fmt.Sprintf("^^<%s#%s>", XMLSchemaNamespace, splits[1])
}
