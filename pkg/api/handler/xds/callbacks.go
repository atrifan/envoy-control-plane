package xds

import (
	"context"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"sync"
	log "github.com/sirupsen/logrus"
)

func (xds *Callbacks) Report() {
	xds.mu.Lock()
	defer xds.mu.Unlock()
	log.WithFields(log.Fields{"fetches": xds.Fetches, "requests": xds.Requests}).Info("xds.Report()  callbacks")
}
func (xds *Callbacks) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	log.Infof("OnStreamOpen %d open for %s", id, typ)
	return nil
}

func (xds *Callbacks) OnStreamClosed(id int64) {
	log.Infof("OnStreamClosed %d closed", id)
}
func (xds *Callbacks) OnStreamRequest(int64, *v2.DiscoveryRequest) error {
	log.Infof("OnStreamRequest")
	xds.mu.Lock()
	defer xds.mu.Unlock()
	xds.Requests++
	if xds.Signal != nil {
		close(xds.Signal)
		xds.Signal = nil
	}

	return nil
}
func (xds *Callbacks) OnStreamResponse(int64, *v2.DiscoveryRequest, *v2.DiscoveryResponse) {
	log.Infof("OnStreamResponse...")
	xds.Report()
}
func (xds *Callbacks) OnFetchRequest(ctx context.Context, req *v2.DiscoveryRequest) error {
	log.Infof("OnFetchRequest...")
	xds.mu.Lock()
	defer xds.mu.Unlock()
	xds.Fetches++
	if xds.Signal != nil {
		close(xds.Signal)
		xds.Signal = nil
	}
	return nil
}
func (xds *Callbacks) OnFetchResponse(*v2.DiscoveryRequest, *v2.DiscoveryResponse) {}

type Callbacks struct {
	Signal   chan struct{}
	Fetches  int
	Requests int
	mu       sync.Mutex
}
