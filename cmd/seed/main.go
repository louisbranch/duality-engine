// Package main provides a CLI for seeding the local development database
// with demo data by exercising the full MCP->game stack, or by generating
// dynamic scenarios directly via gRPC.
package main

import (
	"context"
	"flag"
	"os"

	"github.com/louisbranch/fracturing.space/internal/platform/config"
	"github.com/louisbranch/fracturing.space/internal/tools/cli"

	seedcmd "github.com/louisbranch/fracturing.space/internal/cmd/seed"
)

func main() {
	cfg, err := seedcmd.ParseConfig(flag.CommandLine, os.Args[1:])
	if err != nil {
		config.Exitf("Error: %v", err)
	}

	ctx, stop := cli.WithSignalTimeout(context.Background(), cfg.Timeout)
	defer stop()

	if err := seedcmd.Run(ctx, cfg, os.Stdout, os.Stderr); err != nil {
		config.Exitf("Error: %v", err)
	}
}
