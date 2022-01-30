# This Makefile is an easy way to run common operations.
# Execute commands like this:
# * make

# Load the environment variables.
include .env

.PHONY: default
default: amb

################################################################################
# Setup app
################################################################################

.PHONY: amb
amb:
	go run cmd/amb/main.go

# Update dependencies of other repos using Ambient.
.PHONY: update-children
update-children:
	cd ../plugin && go get github.com/ambientkit/ambient@$(shell git rev-parse HEAD) && go mod tidy -compat=1.17
	cd ../ambient-template && go get github.com/ambientkit/ambient@$(shell git rev-parse HEAD) && go mod tidy -compat=1.17
	cd ../amb && go get github.com/ambientkit/ambient@$(shell git rev-parse HEAD) && go mod tidy -compat=1.17