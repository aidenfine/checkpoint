-include .env

format:
	go fmt ./...

test:
	go test ./...

build:
	go build checkpointmiddleware.go
	
pre-commit:
	$(format)
	$(test)
	$(build)