# This Makefile is an easy way to run common operations.
# Execute commands like this:
# * make
# * make update-children
# * make update-all

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

################################################################################
# Update dependencies
################################################################################

# Run go mod tidy.
.PHONY: tidy
tidy:
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