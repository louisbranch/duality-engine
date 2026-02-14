// Package main provides maintenance utilities.
package main

import (
	"context"
	"flag"
	"os"

	"github.com/louisbranch/fracturing.space/internal/platform/config"
	"github.com/louisbranch/fracturing.space/internal/tools/cli"
	"github.com/louisbranch/fracturing.space/internal/tools/maintenance"
)

func main() {
	cfg, err := maintenance.ParseConfig(flag.CommandLine, os.Args[1:])
	if err != nil {
		config.Exitf("Error: %v", err)
	}

	ctx, stop := cli.WithSignalTimeout(context.Background(), cfg.Timeout)
	defer stop()

	if err := maintenance.Run(ctx, cfg, os.Stdout, os.Stderr); err != nil {
		config.Exitf("Error: %v", err)
	}
}
