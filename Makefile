-include .env

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

run:
	SERVICE_URL=${SERVICE_URL} LOG_LEVEL=${LOG_LEVEL} PORT=${PORT} go run cmd/main.go
