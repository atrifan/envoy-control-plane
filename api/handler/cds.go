package handler

import (
	"context"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
)

type (
	ClusterDiscoveryServiceServer interface {
		StreamClusters(v2.ClusterDiscoveryService_StreamClustersServer) error
		DeltaClusters(v2.ClusterDiscoveryService_DeltaClustersServer) error
		FetchClusters(context.Context, *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error)
	}

	ControlPlane struct {}
)

func (handler *ControlPlane) FetchClusters(context.Context, *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	return &v2.DiscoveryResponse{}, nil
}

func (handler *ControlPlane) StreamClusters(v2.ClusterDiscoveryService_StreamClustersServer) error {
	return nil
}

func (handler *ControlPlane)DeltaClusters(v2.ClusterDiscoveryService_DeltaClustersServer) error {
	return nil
}

func NewControlPlaneCDS() ClusterDiscoveryServiceServer {
	return &ControlPlane{}
}
