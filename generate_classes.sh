#!/usr/bin/env bash
protoc -I /Users/trifan/go/src/github.com/atrifan/envoy-plane/third_party \
    -I /Users/trifan/DEVELOPMENT/envoy/api \
    -I /usr/local/Cellar/go/1.12.5/src/github.com/gogo/googleapis \
    -I /usr/local/Cellar/go/1.12.5/src/github.com/gogo/protobuf \
    -I /usr/local/Cellar/go/1.12.5/src/github.com/envoyproxy/protoc-gen-validate \
    --go_out=plugins=grpc:api \
    /Users/trifan/DEVELOPMENT/envoy/api/envoy/api/v2/cds.proto \
    /Users/trifan/DEVELOPMENT/envoy/api/envoy/api/v2/rds.proto \
    /Users/trifan/DEVELOPMENT/envoy/api/envoy/api/v2/eds.proto

protoc -I /Users/trifan/go/src/github.com/atrifan/envoy-plane/third_party \
    -I /Users/trifan/DEVELOPMENT/envoy/api \
    -I /usr/local/Cellar/go/1.12.5/src/github.com/gogo/googleapis \
    -I /usr/local/Cellar/go/1.12.5/src/github.com/gogo/protobuf \
    -I /usr/local/Cellar/go/1.12.5/src/github.com/envoyproxy/protoc-gen-validate \
    --go_out=. \
    /Users/trifan/DEVELOPMENT/envoy/api/envoy/api/v2/discovery.proto