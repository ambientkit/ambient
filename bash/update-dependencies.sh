#!/bin/bash

# Exit when any command fails.
set -e

AMBIENTHASH=$(git rev-parse HEAD)

cd ../plugin
go get github.com/ambientkit/ambient@${AMBIENTHASH}
go mod tidy -compat=1.17
git add go.mod
git add go.sum
go test ./...
git commit -m "Update dependencies (automated)"
git push
PLUGINHASH=$(git rev-parse HEAD)

cd ../ambient-template
go get github.com/ambientkit/ambient@${AMBIENTHASH}
go get github.com/ambientkit/plugin@${PLUGINHASH}
go mod tidy -compat=1.17
git add go.mod
git add go.sum
go test ./...
git commit -m "Update dependencies (automated)"
git push

cd ../amb
go get github.com/ambientkit/ambient@${AMBIENTHASH}
go get github.com/ambientkit/plugin@${PLUGINHASH}
go mod tidy -compat=1.17
git add go.mod
git add go.sum
go test ./...
git commit -m "Update dependencies (automated)"