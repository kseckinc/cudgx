#!/bin/bash

mkdir -p bin

go fmt ./...
go vet ./...

export GO111MODULE="on"

METRIC_BENCHMARK_NAME="gf.cudgx.sample.benchmark"
go build -o bin/${METRIC_BENCHMARK_NAME} ./sample/benchmark/main.go


