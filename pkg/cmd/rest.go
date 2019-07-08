package cmd

import (
	"context"
	"github.com/atrifan/envoy-plane/pkg/api/handler/grpc"
	"github.com/atrifan/envoy-plane/pkg/api/rest"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/atrifan/envoy-plane/pkg/api/service/v1"
)


// RunServer runs gRPC server and HTTP gateway
func RunServer(ctx context.Context, grpcPort uint, httpPort uint) {

	v1API := v1.NewToDoServiceServer()

	// run HTTP gateway
	go func() {
		_ = grpc.RunServer(ctx, v1API, grpcPort)
	}()

	go func() {
		_ = rest.RunServer(ctx, grpcPort, httpPort)
	}()

	<- ctx.Done()
}