#!/bin/bash

mkdir -p bin

go fmt ./...
go vet ./...

export GO111MODULE="on"

METRIC_CONSUMER_NAME="gf.cudgx.consumer"
go build -o bin/${METRIC_CONSUMER_NAME} ./cmd/consumer/main.go
