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
