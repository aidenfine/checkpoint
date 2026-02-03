-include .env

lint:
	go fmt ./...

test:
	go test ./...

build:
	go build checkpoint.go

benchmark:
	go test -v -bench=. -benchmem

load-test:
	vegeta attack -rate=10000 -duration=30s -targets=targets.txt | vegeta report -type=text > report.txt

high-load-test:
	vegeta attack -rate=25000 -duration=30s -targets=targets.txt | vegeta report -type=text > report.txt
pre-commit:
	$(format)
	$(test)
	$(build)