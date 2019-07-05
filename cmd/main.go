package main

import (
	"context"
	"github.com/atrifan/envoy-plane/pkg/api/handler"
)



func main() {
	ctx := context.Background()
	handler.InitXds(ctx)
}

