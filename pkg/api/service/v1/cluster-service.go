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
	clusterServiceApiVersion = "v1"
)

var Version int32

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type clusterServiceServer struct {
	cache cache.SnapshotCache
}

func (self *clusterServiceServer) ReadAllClustersForNode(ctx context.Context, req *v1.ReadAllRequestForNodeCluster) (*v1.ReadAllResponseForNodeCluster, error) {
	snapshot, err := self.getFromCache(self.cache, req.NodeId)

	response := &v1.ReadAllResponseForNodeCluster{
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

func (self *clusterServiceServer) CreateCluster(ctx context.Context, req *v1.CreateRequestCluster) (*v1.CreateResponseCluster, error) {

	var c []cache.Resource
	atomic.AddInt32(&Version, 1)

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

	snap := cache.NewSnapshot(fmt.Sprint(Version), nil, c, nil, nil)
	err := self.cache.SetSnapshot(req.NodeId, snap)

	return &v1.CreateResponseCluster{
		Version: req.Version,
		NodeId: req.NodeId,
	}, err
}

func (self *clusterServiceServer) DeleteCluster(ctx context.Context, request *v1.DeleteRequestCluster) (*v1.DeleteResponseCluster, error) {
	panic("implement me")

}


func (self *clusterServiceServer) getFromCache(cache cache.SnapshotCache, nodeID string) (cache.Snapshot, error){

	snapshot, err := cache.GetSnapshot(nodeID)

	return snapshot, err
}
// NewToDoServiceServer creates ToDo service
func NewClusterServiceServer(cache cache.SnapshotCache) v1.ClusterServiceServer {
	return &clusterServiceServer{
		cache: cache,
	}
}
