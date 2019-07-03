package main

import (
	"github.com/atrifan/envoy-plane/pkg/api/handler"
	"google.golang.org/grpc"
	"net"
	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
)

func main() {

	var clusters, endpoints, routes, listeners []cache.Resource

	snapshotCache := cache.NewSnapshotCache(false, handler.NewHashFunction(), nil)
	snapshot := cache.NewSnapshot("1.0", endpoints, clusters, routes, listeners)
	_ = snapshotCache.SetSnapshot("node1", snapshot)
	server := xds.NewServer(snapshotCache, handler.NewControlPlaneXDS())
	grpcServer := grpc.NewServer()
	lis, _ := net.Listen("tcp", ":8080")

	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	api.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	api.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	api.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	api.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			// error handling
		}
	}()
}
