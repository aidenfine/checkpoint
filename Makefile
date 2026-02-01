-include .env

lint:
	go fmt ./...

test:
	go test ./...

build:
	go build checkpoint.go
	
pre-commit:
	$(format)
	$(test)
	$(build)