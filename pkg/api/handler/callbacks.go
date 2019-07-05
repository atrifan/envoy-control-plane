package handler

import (
	"context"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"sync"
	log "github.com/sirupsen/logrus"
)

func (handler *Callbacks) Report() {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	log.WithFields(log.Fields{"fetches": handler.fetches, "requests": handler.requests}).Info("handler.Report()  callbacks")
}
func (handler *Callbacks) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	log.Infof("OnStreamOpen %d open for %s", id, typ)
	return nil
}

func (handler *Callbacks) OnStreamClosed(id int64) {
	log.Infof("OnStreamClosed %d closed", id)
}
func (handler *Callbacks) OnStreamRequest(int64, *v2.DiscoveryRequest) error {
	log.Infof("OnStreamRequest")
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.requests++
	if handler.signal != nil {
		close(handler.signal)
		handler.signal = nil
	}

	return nil
}
func (handler *Callbacks) OnStreamResponse(int64, *v2.DiscoveryRequest, *v2.DiscoveryResponse) {
	log.Infof("OnStreamResponse...")
	handler.Report()
}
func (handler *Callbacks) OnFetchRequest(ctx context.Context, req *v2.DiscoveryRequest) error {
	log.Infof("OnFetchRequest...")
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.fetches++
	if handler.signal != nil {
		close(handler.signal)
		handler.signal = nil
	}
	return nil
}
func (handler *Callbacks) OnFetchResponse(*v2.DiscoveryRequest, *v2.DiscoveryResponse) {}

type Callbacks struct {
	signal   chan struct{}
	fetches  int
	requests int
	mu       sync.Mutex
}
