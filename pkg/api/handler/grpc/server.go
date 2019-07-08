package grpc

import (
	"context"
	"fmt"
	v1 "github.com/atrifan/envoy-plane/pkg/api/rest/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port uint) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d",port))
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
