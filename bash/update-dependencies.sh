#!/bin/bash

# This script will update dependencies for the main Ambient repos as long
# as they are all downloaded locally.

# Exit when any command fails.
set -e

AMBIENTHASH=$(git rev-parse HEAD)

echo "Updating repo: plugin"
cd ../plugin
go get github.com/ambientkit/ambient@${AMBIENTHASH}
go mod tidy -compat=1.17
git add go.mod
git add go.sum
go test ./...
git commit -m "Update dependencies (automated)"
git push
PLUGINHASH=$(git rev-parse HEAD)

echo "Updating repo: ambient-template"
cd ../ambient-template
go get github.com/ambientkit/ambient@${AMBIENTHASH}
go get github.com/ambientkit/plugin@${PLUGINHASH}
go mod tidy -compat=1.17
git add go.mod
git add go.sum
go test ./...
git commit -m "Update dependencies (automated)"
git push

echo "Updating repo: amb"
cd ../amb
go get github.com/ambientkit/ambient@${AMBIENTHASH}
go get github.com/ambientkit/plugin@${PLUGINHASH}
go mod tidy -compat=1.17
git add go.mod
git add go.sum
go test ./...
git commit -m "Update dependencies (automated)"
git push