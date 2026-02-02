-include .env

lint:
	go fmt ./...

test:
	go test ./...

build:
	go build checkpoint.go

benchmark:
	go test -v -bench=. -benchmem
	
pre-commit:
	$(format)
	$(test)
	$(build)