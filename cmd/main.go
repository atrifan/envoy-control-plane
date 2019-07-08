package main

import (
	"context"
	"github.com/atrifan/envoy-plane/pkg/cmd"
)



func main() {
	ctx := context.Background()
	cmd.InitXds(ctx)
}

