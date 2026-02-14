// Package main provides a CLI for running Lua scenario scripts.
package main

import (
	"context"
	"flag"
	"os"

	"github.com/louisbranch/fracturing.space/internal/platform/config"
	"github.com/louisbranch/fracturing.space/internal/tools/cli"

	scenariocmd "github.com/louisbranch/fracturing.space/internal/cmd/scenario"
)

func main() {
	cfg, err := scenariocmd.ParseConfig(flag.CommandLine, os.Args[1:])
	if err != nil {
		config.Exitf("Error: %v", err)
	}

	ctx, stop := cli.WithSignalContext(context.Background())
	defer stop()

	if err := scenariocmd.Run(ctx, cfg, os.Stdout, os.Stderr); err != nil {
		config.Exitf("Error: %v", err)
	}
}
