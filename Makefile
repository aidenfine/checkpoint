format:
	go fmt ./...

test:
	go test ./...

build:
	go build cmd/main.go
	
pre-commit:
	$(format)
	$(test)
	$(build)
