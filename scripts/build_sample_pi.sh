#!/bin/bash

mkdir -p bin

go fmt ./...
go vet ./...

export GO111MODULE="on"

METRIC_PI_NAME="gf.cudgx.sample.pi"
go build -o bin/${METRIC_PI_NAME} ./sample/pi/main.go

