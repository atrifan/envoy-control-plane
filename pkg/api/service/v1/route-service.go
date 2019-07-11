package v1

import (
	"context"
	"fmt"
	v1 "github.com/atrifan/envoy-plane/pkg/api/handler/rest/v1"
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"sync/atomic"
)

const (
	// apiVersion is version of API is provided by server
	routeServiceApiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type routeServiceServer struct {
	cache cache.SnapshotCache
}

func (self *routeServiceServer) ReadAllRoutesForNodeAndRoute(ctx context.Context, req *v1.ReadResponseForNodeRoute) (*v1.Route, error) {
	snapshot, err := GetFromCache(self.cache, req.NodeId)

	response := &v1.Route{}

	for _, value := range snapshot.Routes.Items {
		switch value := value.(type) {
		case *v2.RouteConfiguration:
			routeInstance := &v1.Route{
				Name: value.Name,
				Domains: value.VirtualHosts[0].Domains,
			}
			return routeInstance, nil
		}
	}

	return response, err
}

func (self *routeServiceServer) CreateRoute(ctx context.Context, req *v1.CreateRequestRoute) (*v1.CreateResponseRoute, error) {

	var r []cache.Resource
	atomic.AddInt32(&Version, 1)

	var routes []route.Route

	for _, entry := range req.Route.RouteInfo {
		switch entry.Type {
			case "prefix":
				routes = append(routes, route.Route{
					Match: route.RouteMatch {
						PathSpecifier: &route.RouteMatch_Prefix{
							Prefix: entry.Value,
						},
					},
					Action: &route.Route_Route{
						Route: &route.RouteAction{
							ClusterSpecifier: &route.RouteAction_Cluster{
								Cluster: entry.Route.Cluster,
							},
						},
					},
				})
				break
			case "path":
				routes = append(routes, route.Route{
					Match: route.RouteMatch{
						PathSpecifier: &route.RouteMatch_Path{
							Path: entry.Value,
						},
					},
					Action: &route.Route_Route{
						Route: &route.RouteAction{
							ClusterSpecifier: &route.RouteAction_Cluster{
								Cluster: entry.Route.Cluster,
							},
						},
					},

				})
				break
		}

	}
	routeEntry := &v2.RouteConfiguration{
		Name: req.Route.Name,
		VirtualHosts: []route.VirtualHost{{
			Name: "local_service",
			Domains: req.Route.Domains,
			Routes: routes,
		}},
	}

	r = append(r, routeEntry)

	oldClusters, _ := GetOldClusters(self.cache, req.NodeId)
	oldListeners, _ := GetOldListeners(self.cache, req.NodeId)

	snap := cache.NewSnapshot(fmt.Sprint(Version), nil, oldClusters, r, oldListeners)
	err := self.cache.SetSnapshot(req.NodeId, snap)

	return &v1.CreateResponseRoute{
		Version: fmt.Sprint(Version),
		NodeId: req.NodeId,
	}, err
}


// NewToDoServiceServer creates ToDo service
func NewRouteServiceServer(cache cache.SnapshotCache) v1.RouteServiceServer {
	return &routeServiceServer{
		cache: cache,
	}
}
