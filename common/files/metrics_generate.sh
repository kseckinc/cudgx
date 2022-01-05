#/bin/bash

protoc --proto_path=files  --go_out=mod --go_opt=paths=source_relative files/metrics.proto
protoc --proto_path=files  --go_out=mod --go_opt=paths=source_relative files/streaming.proto

