package main

import (
	"fmt"
	"github.com/atrifan/envoy-plane/api/handler"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"log"
	"net"
	"google.golang.org/grpc"
)
// main start a gRPC server and waits for connection
func main() {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	var cds handler.ClusterDiscoveryServiceServer
	cds = handler.NewControlPlaneCDS()
	// start the server
	v2.RegisterClusterDiscoveryServiceServer(grpcServer, cds)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
