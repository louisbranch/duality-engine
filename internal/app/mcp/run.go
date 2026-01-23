package mcp

import (
	"context"

	"github.com/louisbranch/duality-engine/internal/mcp/service"
)

// Run starts the MCP app with the provided gRPC address.
func Run(ctx context.Context, addr string) error {
	return service.Run(ctx, service.Config{
		GRPCAddr:  addr,
		Transport: service.TransportStdio,
	})
}
