#!/bin/bash

./scripts/build.sh

docker build -f docker/api.Dockerfile .  -t 172.16.16.172:12380/cudgx/api:v0.2
docker push 172.16.16.172:12380/cudgx/api:v0.2

docker build -f docker/gateway.Dockerfile .  -t 172.16.16.172:12380/cudgx/gateway:v0.2
docker push 172.16.16.172:12380/cudgx/gateway:v0.2

docker build -f docker/consumer.Dockerfile  .  -t 172.16.16.172:12380/cudgx/consumer:v0.2
docker push 172.16.16.172:12380/cudgx/consumer:v0.2

docker build -f docker/pi.Dockerfile . -t 172.16.16.172:12380/cudgx/sample-pi:v0.2
docker push 172.16.16.172:12380/cudgx/sample-pi:v0.2

docker build -f docker/benchmark.Dockerfile . -t 172.16.16.172:12380/cudgx/sample-benchmark:v0.2
docker push 172.16.16.172:12380/cudgx/sample-benchmark:v0.2


