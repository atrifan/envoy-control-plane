package main

import (
	"context"
	"github.com/atrifan/envoy-plane/cmd/xds"
)



func main() {
	ctx := context.Background()
	xds.InitXds(ctx)
}

