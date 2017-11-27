generate:
	go generate gen/aws/generators/main.go
build: generate
	go build
