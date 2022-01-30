#!/bin/bash

AMBIENTHASH=$(git rev-parse HEAD)

cd ../plugin
PLUGINHASH=$(git rev-parse HEAD)
git add go.mod
git add go.sum
# git commit -m "Update dependencies"
go get github.com/ambientkit/ambient@${AMBIENTHASH}}
go mod tidy -compat=1.17

cd ../ambient-template
git add go.mod
git add go.sum
# git commit -m "Update dependencies"
go get github.com/ambientkit/ambient@${AMBIENTHASH}}
go mod tidy -compat=1.17

cd ../amb
git add go.mod
git add go.sum
# git commit -m "Update dependencies"
go get github.com/ambientkit/ambient@${AMBIENTHASH}}
go mod tidy -compat=1.17