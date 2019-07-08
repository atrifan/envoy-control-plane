protoc --proto_path=api/rest/v1 \
    --proto_path=vendor/github.com/grpc-ecosystem/grpc-gateway \
    --proto_path=vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --go_out=plugins=grpc:pkg/api/rest/v1 cluster-service.proto
protoc --proto_path=api/rest/v1 \
    --proto_path=vendor/github.com/grpc-ecosystem/grpc-gateway \
    --proto_path=vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=logtostderr=true:pkg/api/rest/v1 cluster-service.proto
protoc --proto_path=api/rest/v1 \
    --proto_path=vendor/github.com/grpc-ecosystem/grpc-gateway \
    --proto_path=vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:api/swagger/rest/v1 cluster-service.proto