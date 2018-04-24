//go:generate go run $GOFILE properties.go paramsdoc.go mocks.go fetchers.go services.go commands.go acceptance_mocks.go

//go:generate gofmt -s -w ../../../aws
//go:generate goimports -w ../../../aws

//go:generate gofmt -s -w ../../../aws/services
//go:generate goimports -w ../../../aws/services

//go:generate gofmt -s -w ../../../aws/fetch
//go:generate goimports -w ../../../aws/fetch

//go:generate gofmt -s -w ../../../cloud/properties
//go:generate goimports -w ../../../cloud/properties

//go:generate gofmt -s -w ../../../cloud/rdf
//go:generate goimports -w ../../../cloud/rdf

//go:generate gofmt -s -w ../../../aws/spec
//go:generate goimports -w ../../../aws/spec

//go:generate gofmt -s -w ../../../acceptance/aws
//go:generate goimports -w ../../../acceptance/aws

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	"text/template"
)

var (
	ROOT_DIR = filepath.Join("..", "..", "..")

	FETCHERS_DIR         = filepath.Join(ROOT_DIR, "aws", "fetch")
	SERVICES_DIR         = filepath.Join(ROOT_DIR, "aws", "services")
	SPEC_DIR             = filepath.Join(ROOT_DIR, "aws", "spec")
	AWSAT_DIR            = filepath.Join(ROOT_DIR, "acceptance", "aws")
	DOC_DIR              = filepath.Join(ROOT_DIR, "aws", "doc")
	CLOUD_PROPERTIES_DIR = filepath.Join(ROOT_DIR, "cloud", "properties")
	CLOUD_RDF_DIR        = filepath.Join(ROOT_DIR, "cloud", "rdf")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("[+] ")
	flag.Parse()

	// fetchers
	generateFetcherFuncs()
	generateServicesFuncs()

	// mocks
	generateTestMocks()

	// commands
	generateCommands()
	generateAcceptanceMocks()
	generateAcceptanceFactory()

	// properties
	generateProperties()
	generateRDFProperties()

	// doc
	generateParamsDocLookup()
}

func writeTemplateToFile(templ *template.Template, data interface{}, dir, filename string) {
	var buff bytes.Buffer
	if err := templ.Execute(&buff, data); err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dir, filename)
	if err := ioutil.WriteFile(path, buff.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}

	log.Printf("generated %s", relativePathToRoot(path))
}

func relativePathToRoot(path string) string {
	rel, _ := filepath.Rel(ROOT_DIR, path)
	return rel
}
