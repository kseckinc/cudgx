#!/bin/bash

mkdir -p bin

go fmt ./...
go vet ./...

export GO111MODULE="on"

METRIC_API_NAME="gf.cudgx.api"
go build -o bin/${METRIC_API_NAME} ./cmd/api/main.go

METRIC_GATEWAY_NAME="gf.cudgx.gateway"
go build -o bin/${METRIC_GATEWAY_NAME} ./cmd/gateway/main.go

METRIC_CONSUMER_NAME="gf.cudgx.consumer"
go build -o bin/${METRIC_CONSUMER_NAME} ./cmd/consumer/main.go

METRIC_PI_NAME="gf.cudgx.sample.pi"
go build -o bin/${METRIC_PI_NAME} ./sample/pi/main.go

METRIC_BENCHMARK_NAME="gf.cudgx.sample.benchmark"
go build -o bin/${METRIC_BENCHMARK_NAME} ./sample/benchmark/main.go


