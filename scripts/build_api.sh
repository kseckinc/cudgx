#!/bin/bash

mkdir -p bin

go fmt ./...
go vet ./...

export GO111MODULE="on"

METRIC_API_NAME="gf.cudgx.api"
go build -o bin/${METRIC_API_NAME} ./cmd/api/main.go

