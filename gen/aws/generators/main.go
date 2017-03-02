//go:generate go run $GOFILE drivers.go fetchers.go
//go:generate gofmt -s -w ../../../aws
//go:generate goimports -w ../../../aws
//go:generate gofmt -s -w ../../../aws/driver
//go:generate goimports -w ../../../aws/driver

package main

import "path/filepath"

var (
	ROOT_DIR = filepath.Join("..", "..", "..")

	FETCHERS_DIR = filepath.Join(ROOT_DIR, "aws")
	DRIVERS_DIR  = filepath.Join(ROOT_DIR, "aws", "driver")
)

func main() {
	// fetchers
	generateFetcherFuncs()

	// drivers, templates
	generateDriverFuncs()
	generateTemplateTemplates()
	generateDriverTypes()
}
