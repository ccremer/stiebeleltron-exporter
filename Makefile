.PHONY: build fmt dist clean test run
SHELL := /usr/bin/env bash

build: fmt
	@go build .

fmt:
	@[[ -z $$(go fmt ./...) ]]

dist: fmt
	@goreleaser release --snapshot --rm-dist --skip-sign

clean:
	@rm -rf stiebeleltron-exporter c.out dist

test:
	@go test -coverprofile c.out ./...

run:
	@go run . -v
