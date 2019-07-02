package handler

import (
	"context"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
)

type (
	Callbacks interface {
		// OnStreamOpen is called once an xDS stream is open with a stream ID and the type URL (or "" for ADS).
		// Returning an error will end processing and close the stream. OnStreamClosed will still be called.
		OnStreamOpen(context.Context, int64, string) error
		// OnStreamClosed is called immediately prior to closing an xDS stream with a stream ID.
		OnStreamClosed(int64)
		// OnStreamRequest is called once a request is received on a stream.
		// Returning an error will end processing and close the stream. OnStreamClosed will still be called.
		OnStreamRequest(int64, *v2.DiscoveryRequest) error
		// OnStreamResponse is called immediately prior to sending a response on a stream.
		OnStreamResponse(int64, *v2.DiscoveryRequest, *v2.DiscoveryResponse)
		// OnFetchRequest is called for each Fetch request. Returning an error will end processing of the
		// request and respond with an error.
		OnFetchRequest(context.Context, *v2.DiscoveryRequest) error
		// OnFetchResponse is called immediately prior to sending a response.
		OnFetchResponse(*v2.DiscoveryRequest, *v2.DiscoveryResponse)
	}

	NodeHash interface {
		ID(node *core.Node) string
	}

	ControlPlane struct {}

	ControlPlaneHash struct {}
)

func (handler *ControlPlane) OnStreamOpen(context.Context, int64, string) error {
	return nil
}

func (handler *ControlPlane) OnStreamClosed(int64) {
	return
}

func (handler *ControlPlane) OnStreamRequest(int64, *v2.DiscoveryRequest) error {
	return nil
}

func (handler *ControlPlane) OnFetchRequest(context.Context, *v2.DiscoveryRequest) error {
	return nil
}

func (handler *ControlPlane) OnStreamResponse(int64, *v2.DiscoveryRequest, *v2.DiscoveryResponse) {
	return
}

func (handler *ControlPlane) OnFetchResponse(*v2.DiscoveryRequest, *v2.DiscoveryResponse) {
	return
}

func (handler *ControlPlaneHash) ID(node *core.Node) string {
	return "someHash"
}

func NewControlPlaneXDS() Callbacks {
	return &ControlPlane{}
}

func NewHashFunction() NodeHash {
	return &ControlPlaneHash{}
}
