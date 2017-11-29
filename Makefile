test:
	@echo Running tests
	@go test ./...

generate:
	@echo Generating boilerplate code
	@go generate gen/aws/generators/main.go

build: generate test
	@echo Building application binary
	@go build
