

include Makefile.common

.PHONY: all clean test gotest docker docker-gateway docker-api docker-consumer docker-pi docker-benchmark docker-push push-gateway push-consumer push-api push-pi push-benchmark

default: all buildsucc


buildsucc:
	@echo Build Cudgx successfully!

all: dev gateway consumer pi benchmark

dev: check
	@>&2 echo "Great, all tests passed."

check: fmt vet

fmt:
	@echo "gofmt (simplify)"
	@gofmt -s -l -w $(FILES) 2>&1 | $(FAIL_ON_STDOUT)

vet:
	@echo "vet "
	@echo  $(PACKAGES_CUDGX_TESTS)
	$(GO) vet -all $(PACKAGES_CUDGX_TESTS) 2>&1 | $(FAIL_ON_STDOUT)

gateway:
	CGO_ENABLED=0 $(GOBUILD) $(RACE_FLAG) -ldflags '$(LDFLAGS) $(CHECK_FLAG)' -o bin/gf.cudgx.gateway ./cmd/gateway/main.go

consumer:
	CGO_ENABLED=0 $(GOBUILD) $(RACE_FLAG) -ldflags '$(LDFLAGS) $(CHECK_FLAG)' -o bin/gf.cudgx.consumer ./cmd/consumer/main.go

api:
	CGO_ENABLED=0 $(GOBUILD) $(RACE_FLAG) -ldflags '$(LDFLAGS) $(CHECK_FLAG)' -o bin/gf.cudgx.api ./cmd/api/main.go

pi:
	CGO_ENABLED=0 $(GOBUILD) $(RACE_FLAG) -ldflags '$(LDFLAGS) $(CHECK_FLAG)' -o bin/gf.cudgx.sample.pi ./sample/pi/main.go

benchmark:
	CGO_ENABLED=0 $(GOBUILD) $(RACE_FLAG) -ldflags '$(LDFLAGS) $(CHECK_FLAG)' -o bin/gf.cudgx.sample.benchmark ./sample/benchmark/main.go

docker: docker-gateway docker-consumer docker-pi docker-benchmark docker-api buildsucc

docker-gateway: gateway
	@docker build -f docker/gateway.Dockerfile .  -t 172.16.16.172:12380/cudgx/gateway:$(IMAGE_VERSION)

docker-consumer: consumer
	@docker build -f docker/consumer.Dockerfile  .  -t 172.16.16.172:12380/cudgx/consumer:$(IMAGE_VERSION)

docker-api: api
	@docker build -f docker/api.Dockerfile  .  -t 172.16.16.172:12380/cudgx/api:$(IMAGE_VERSION)

docker-pi: pi
	@docker build -f docker/pi.Dockerfile . -t 172.16.16.172:12380/cudgx/sample-pi:$(IMAGE_VERSION)

docker-benchmark: benchmark
	@docker build -f docker/benchmark.Dockerfile . -t 172.16.16.172:12380/cudgx/sample-benchmark:$(IMAGE_VERSION)


docker-push: docker push-gateway push-consumer push-pi push-api push-benchmark


push-gateway: docker-gateway
	docker push 172.16.16.172:12380/cudgx/gateway:$(IMAGE_VERSION)

push-consumer: docker-consumer
	docker push 172.16.16.172:12380/cudgx/consumer:$(IMAGE_VERSION)

push-api: docker-api
	docker push 172.16.16.172:12380/cudgx/api:$(IMAGE_VERSION)

push-pi: docker-pi
	docker push 172.16.16.172:12380/cudgx/sample-pi:$(IMAGE_VERSION)

push-benchmark: docker-benchmark
	docker push 172.16.16.172:12380/cudgx/sample-benchmark:$(IMAGE_VERSION)

# Quick start
# Pull images from dockerhub and run
docker-run-linux:
	sh ./run-for-linux.sh

docker-run-mac:
	sh ./run-for-mac.sh

docker-container-stop:
	docker ps -aq | xargs docker stop
	docker ps -aq | xargs docker rm

docker-image-rm:
	docker image prune --force --all

# Immersive experience
# Compile and run by docker-compose
docker-compose-start:
	docker-compose up -d

docker-compose-stop:
	docker-compose down

docker-compose-build:
	docker-compose build

#USE make TARGET version=xx override version
version ?= latest

docker-tag:
	docker tag cudgx_api:latest galaxyfuture/cudgx-api:${version}
	docker tag cudgx_gateway:latest galaxyfuture/cudgx-gateway:${version}
	docker tag cudgx_consumer:latest galaxyfuture/cudgx-consumer:${version}
	docker tag cudgx_sample_pi:latest galaxyfuture/cudgx-sample-pi:${version}
	docker tag cudgx_sample_benchmark:latest galaxyfuture/cudgx-sample-benchmark:${version}

docker-push-hub:
	docker push galaxyfuture/cudgx-api:${version}
	docker push galaxyfuture/cudgx-gateway:${version}
	docker push galaxyfuture/cudgx-consumer:${version}
	docker push galaxyfuture/cudgx-sample-pi:${version}
	docker push galaxyfuture/cudgx-sample-benchmark:${version}

docker-hub-all: docker-compose-build docker-tag docker-push-hub






