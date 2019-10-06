
GOBIN=$(shell pwd)/bin

alien:
	GOBIN=$(GOBIN) go install cmd/alieninvasion.go

test:
	go test -cover ./...

nice:
	go fmt ./...

.PHONY: alien test nice