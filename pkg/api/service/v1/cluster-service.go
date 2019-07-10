package v1

import (
	"context"
	"fmt"
	v1 "github.com/atrifan/envoy-plane/pkg/api/handler/rest/v1"
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"sync/atomic"
	"time"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

var version int32

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type clusterServiceServer struct {
	cache cache.SnapshotCache
}

func (self *clusterServiceServer) ReadAllClustersForNode(ctx context.Context, req *v1.ReadAllRequestForNode) (*v1.ReadAllResponseForNode, error) {
	snapshot, err := getFromCache(self.cache, req.NodeId)

	response := &v1.ReadAllResponseForNode{
		NodeId: req.NodeId,
	}

	for _, value := range snapshot.Clusters.Items {
		switch value := value.(type) {
			case *v2.Cluster:
				cluster := &v1.Cluster{
					ClusterName: value.Name,
					LbPolicy: value.LbPolicy.String(),
				}
				for _, host := range value.Hosts {
					cluster.Hosts = append(cluster.Hosts, &v1.Hosts{
						Port: host.GetSocketAddress().GetPortValue(),
						Ip:   host.GetSocketAddress().GetAddress(),
					})
				}
				response.Clusters = append(response.Clusters, cluster)
		}
	}

	return response, err
}

func (self *clusterServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {

	var c []cache.Resource
	atomic.AddInt32(&version, 1)

	for _, cluster := range req.Cluster {
		clusterEntry := &v2.Cluster{
			Name: cluster.ClusterName,
			ConnectTimeout:  2 * time.Second,
			ClusterDiscoveryType: &v2.Cluster_Type{v2.Cluster_STATIC},
			DnsLookupFamily: v2.Cluster_V4_ONLY,
			LbPolicy:        v2.Cluster_ROUND_ROBIN,
		}

		for _, host := range cluster.Hosts {
			clusterEntry.Hosts = append(clusterEntry.Hosts, &core.Address{Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Address:  host.Ip,
					Protocol: core.TCP,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: host.Port,
					},
				},
			}})
		}
		c = append(c, clusterEntry)
	}

	snap := cache.NewSnapshot(fmt.Sprint(version), nil, c, nil, nil)
	err := self.cache.SetSnapshot(req.NodeId, snap)

	return &v1.CreateResponse{
		Version: req.Version,
		NodeId: req.NodeId,
	}, err
}

func (self *clusterServiceServer) Read(context.Context, *v1.ReadRequest) (*v1.ReadResponse, error) {
	panic("implement me")
}

func (self *clusterServiceServer) Update(context.Context, *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	panic("implement me")
}

func (self *clusterServiceServer) Delete(ctx context.Context, request *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	panic("implement me")

}


func getFromCache(cache cache.SnapshotCache, nodeID string) (cache.Snapshot, error){

	snapshot, err := cache.GetSnapshot(nodeID)

	return snapshot, err
}
// NewToDoServiceServer creates ToDo service
func NewToDoServiceServer(cache cache.SnapshotCache) v1.ClusterServiceServer {
	return &clusterServiceServer{
		cache: cache,
	}
}
