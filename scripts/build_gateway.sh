#!/bin/bash

mkdir -p bin

go fmt ./...
go vet ./...

export GO111MODULE="on"

METRIC_GATEWAY_NAME="gf.cudgx.gateway"
go build -o bin/${METRIC_GATEWAY_NAME} ./cmd/gateway/main.go
