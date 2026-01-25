package mcp

import (
	"context"

	"github.com/louisbranch/duality-engine/internal/mcp/service"
)

// Run starts the MCP app with the provided gRPC address, HTTP address, and transport type.
func Run(ctx context.Context, grpcAddr, httpAddr, transport string) error {
	transportKind := service.TransportStdio
	if transport == "http" {
		transportKind = service.TransportHTTP
	}
	
	return service.Run(ctx, service.Config{
		GRPCAddr:  grpcAddr,
		HTTPAddr:  httpAddr,
		Transport: transportKind,
	})
}
