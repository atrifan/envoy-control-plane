package v1

import "github.com/envoyproxy/go-control-plane/pkg/cache"

func GetFromCache(cache cache.SnapshotCache, nodeID string) (cache.Snapshot, error){
	snapshot, err := cache.GetSnapshot(nodeID)
	return snapshot, err
}

func GetOldClusters(cacheInstance cache.SnapshotCache, nodeId string) ([]cache.Resource, error){
	snapshot, err := GetFromCache(cacheInstance, nodeId)

	var oldClusters []cache.Resource

	for _, value := range snapshot.GetResources(cache.ClusterType) {
		oldClusters = append(oldClusters, value)
	}

	return oldClusters, err
}

func GetOldListeners(cacheInstance cache.SnapshotCache, nodeId string) ([]cache.Resource, error) {
	snapshot, err := GetFromCache(cacheInstance, nodeId)
	var oldListeners []cache.Resource
	for _, value := range snapshot.GetResources(cache.ListenerType) {
		oldListeners = append(oldListeners, value)
	}
	return oldListeners, err
}

func GetOldRoutes(cacheInstance cache.SnapshotCache, nodeId string) ([]cache.Resource, error) {
	snapshot, err := GetFromCache(cacheInstance, nodeId)
	var oldRoutes []cache.Resource

	for _, value := range snapshot.GetResources(cache.RouteType) {
		oldRoutes = append(oldRoutes, value)
	}

	return oldRoutes, err
}
