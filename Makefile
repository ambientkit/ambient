# This Makefile is an easy way to run common operations.
# Execute commands like this:
# * make
# * make install
# * make generate
# * make update-children
# * make update-all

# Load the environment variables.
include .env

.PHONY: default
default: generate

################################################################################
# Common
################################################################################

# Install dependencies.
.PHONY: install
install: protoc-install tidy

################################################################################
# Update dependencies
################################################################################

# Update dependencies of other repos using Ambient.
.PHONY: update-children
update-children:
	cd ../plugin && go get github.com/ambientkit/ambient@$(shell cd ../ambient && git rev-parse HEAD) && go mod tidy -compat=1.17
	cd ../ambient-template && go get github.com/ambientkit/ambient@$(shell cd ../ambient && git rev-parse HEAD) && go mod tidy -compat=1.17
	cd ../ambient-template && go get github.com/ambientkit/plugin@$(shell cd ../plugin && git rev-parse HEAD) && go mod tidy -compat=1.17
	cd ../amb && go get github.com/ambientkit/ambient@$(shell cd ../ambient && git rev-parse HEAD) && go mod tidy -compat=1.17
	cd ../amb && go get github.com/ambientkit/plugin@$(shell cd ../plugin && git rev-parse HEAD) && go mod tidy -compat=1.17

# Update all Go dependencies.
.PHONY: update-all
update-all: update-all-go tidy

# Update all Go dependencies.
.PHONY: update-all-go
update-all-go:
	go get -u -f -d ./...

# Run go mod tidy.
.PHONY: tidy
tidy:
	go mod tidy -compat=1.17

################################################################################
# gRPC
################################################################################

# Install protoc to project bin folder to allow generating a Go file from proto file.
.PHONY: protoc-install
protoc-install:
	mkdir -p ./bin
	curl -s -o protoc.zip -L https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-osx-x86_64.zip
	unzip -q protoc.zip -d tempdir
	rm protoc.zip
	cp tempdir/bin/protoc ./bin/
	cp -r tempdir/include ./bin/
	rm -r tempdir
	GOBIN=$(shell pwd)/bin go install github.com/golang/protobuf/protoc-gen-go@latest

# Generate the Go code, grpc code, and tidy.
.PHONY: generate
generate:
	go generate ./...
	@PATH="${PATH}:$(shell pwd)/bin" && protoc -I pkg/grpcp/protobuf/ pkg/grpcp/protobuf/*.proto --go_out=plugins=grpc:pkg/grpcp/protodef/
	go mod tidy -compat=1.17