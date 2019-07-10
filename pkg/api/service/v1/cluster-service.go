package v1

import (
	"context"
	v1 "github.com/atrifan/envoy-plane/pkg/api/handler/rest/v1"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type clusterServiceServer struct {
	cache *cache.SnapshotCache
}

func (v1 *clusterServiceServer) ReadAllClustersForNode(context.Context, *v1.ReadAllRequestForNode) (*v1.ReadAllResponseForNode, error) {
	panic("implement me")
}

func (v1 *clusterServiceServer) Create(context.Context, *v1.CreateRequest) (*v1.CreateResponse, error) {
	panic("implement me")
}

func (v1 *clusterServiceServer) Read(context.Context, *v1.ReadRequest) (*v1.ReadResponse, error) {
	panic("implement me")
}

func (v1 *clusterServiceServer) Update(context.Context, *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	panic("implement me")
}

func (v1 *clusterServiceServer) Delete(context.Context, *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	panic("implement me")
}

// NewToDoServiceServer creates ToDo service
func NewToDoServiceServer(cache *cache.SnapshotCache) v1.ClusterServiceServer {
	return &clusterServiceServer{
		cache: cache,
	}
}
