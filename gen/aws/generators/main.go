//go:generate go run $GOFILE drivers.go fetchers.go
//go:generate gofmt -s -w ../../../cloud/aws/gen_api.go
//go:generate gofmt -s -w ../../../template/driver/aws/gen_template_defs.go
//go:generate gofmt -s -w ../../../template/driver/aws/gen_driver_funcs.go
//go:generate gofmt -s -w ../../../template/driver/aws/gen_drivers.go

package main

import "path/filepath"

var (
	ROOT_DIR = filepath.Join("..", "..", "..")

	FETCHERS_DIR = filepath.Join(ROOT_DIR, "cloud", "aws")
	DRIVERS_DIR  = filepath.Join(ROOT_DIR, "template", "driver", "aws")
)

func main() {
	// fetchers
	generateFetcherFuncs()

	// drivers, templates
	generateDriverFuncs()
	generateTemplateTemplates()
	generateDriverTypes()
}
