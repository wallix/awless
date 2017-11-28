test:
	@echo Running tests
	@go test $$(go list ./... | grep -v /vendor/)

generate:
	@echo Generating boilerplate code
	@go generate gen/aws/generators/main.go

build: generate
	@echo Building application binary
	@go build
