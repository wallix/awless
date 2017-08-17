package triplestore

import (
	"fmt"
	"strings"
)

type XsdType string

var (
	XsdString   = XsdType("xsd:string")
	XsdBoolean  = XsdType("xsd:boolean")
	XsdDateTime = XsdType("xsd:dateTime")

	// 64-bit floating point numbers
	XsdDouble = XsdType("xsd:double")
	// 32-bit floating point numbers
	XsdFloat = XsdType("xsd:float")

	// signed 32 or 64 bit
	XsdInteger = XsdType("xsd:integer")
	// signed (8 bit)
	XsdByte = XsdType("xsd:byte")
	// signed (16 bit)
	XsdShort = XsdType("xsd:short")

	// unsigned 32 or 64 bit
	XsdUinteger = XsdType("xsd:unsignedInt")
	// unsigned 8 bit
	XsdUnsignedByte = XsdType("xsd:unsignedByte")
	// unsigned 16 bit
	XsdUnsignedShort = XsdType("xsd:unsignedShort")
)

const XMLSchemaNamespace = "http://www.w3.org/2001/XMLSchema"

func (x XsdType) NTriplesNamespaced() string {
	splits := strings.Split(string(x), ":")
	if len(splits) != 2 {
		return string(x)
	}

	return fmt.Sprintf("%s#%s", XMLSchemaNamespace, splits[1])
}
