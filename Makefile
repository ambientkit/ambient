# This Makefile is an easy way to run common operations.
# Execute commands like this:
# * make
# * make update-children
# * make update-all

# Load the environment variables.
include .env

.PHONY: default
default: start

################################################################################
# Update dependencies
################################################################################

# Run go mod tidy.
.PHONY: tidy
tidy:
	go generate ./...
	go mod tidy -compat=1.17

# Update dependencies of other repos using Ambient.
.PHONY: update-children
update-children:
	cd ../plugin && go get github.com/ambientkit/ambient@$(shell git rev-parse HEAD) && go mod tidy -compat=1.17
	cd ../ambient-template && go get github.com/ambientkit/ambient@$(shell git rev-parse HEAD) && go mod tidy -compat=1.17
	cd ../amb && go get github.com/ambientkit/ambient@$(shell git rev-parse HEAD) && go mod tidy -compat=1.17

# Update all other repos with latest Ambient dependencies.
.PHONY: update-all
update-all:
	 ./bash/update-dependencies.sh


################################################################################
# gRPC
################################################################################

# Install protoc to project bin folder to allow generating a Go file from proto file.
.PHONY: protoc-install
protoc-install:
	mkdir ./bin
	curl -s -o protoc.zip -L https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-osx-x86_64.zip
	unzip -q protoc.zip -d tempdir
	rm protoc.zip
	cp tempdir/bin/protoc ./bin/
	cp -r tempdir/include ./bin/
	rm -r tempdir
	GOBIN=$(shell pwd)/bin go install github.com/golang/protobuf/protoc-gen-go@latest

# Generate the grpc code.
.PHONY: protoc
protoc:
	@PATH="${PATH}:$(shell pwd)/bin" && protoc -I pkg/grpcp/protobuf/ pkg/grpcp/protobuf/*.proto --go_out=plugins=grpc:pkg/grpcp/protodef/

# Start the build and run process.
.PHONY: start
start: protoc
	@cd cmd/plugin/hello/cmd/plugin && go build -o hello
	go run cmd/server/main.go