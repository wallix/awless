//go:generate go run $GOFILE drivers.go properties.go paramsdoc.go mocks.go fetchers.go services.go
//go:generate gofmt -s -w ../../../aws
//go:generate goimports -w ../../../aws

//go:generate gofmt -s -w ../../../aws/services
//go:generate goimports -w ../../../aws/services

//go:generate gofmt -s -w ../../../aws/fetch
//go:generate goimports -w ../../../aws/fetch

//go:generate gofmt -s -w ../../../aws/driver
//go:generate goimports -w ../../../aws/driver

//go:generate gofmt -s -w ../../../cloud/properties
//go:generate goimports -w ../../../cloud/properties

//go:generate gofmt -s -w ../../../cloud/rdf
//go:generate goimports -w ../../../cloud/rdf

package main

import (
	"flag"
	"path/filepath"
)

var (
	ROOT_DIR = filepath.Join("..", "..", "..")

	FETCHERS_DIR         = filepath.Join(ROOT_DIR, "aws", "fetch")
	SERVICES_DIR         = filepath.Join(ROOT_DIR, "aws", "services")
	DRIVERS_DIR          = filepath.Join(ROOT_DIR, "aws", "driver")
	DOC_DIR              = filepath.Join(ROOT_DIR, "aws", "doc")
	CLOUD_PROPERTIES_DIR = filepath.Join(ROOT_DIR, "cloud", "properties")
	CLOUD_RDF_DIR        = filepath.Join(ROOT_DIR, "cloud", "rdf")
)

func main() {
	flag.Parse()

	// fetchers
	generateFetcherFuncs()
	generateServicesFuncs()

	// mocks
	generateTestMocks()

	// drivers, templates
	generateDriverFuncs()
	generateTemplateTemplates()
	generateDriverTypes()

	// properties
	generateProperties()
	generateRDFProperties()

	// doc
	if true {
		generateParamsDocLookup()
	}
}
