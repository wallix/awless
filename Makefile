test:
	@echo Running tests (with -race flag on) 
	@go test ./... -race

generate:
	@echo Generating commands code: runtime, doc, etc.
	@go generate gen/aws/generators/main.go

build: generate test
	@echo Building application binary
	@go build
