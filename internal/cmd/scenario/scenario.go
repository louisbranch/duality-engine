package scenario

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"time"

	"github.com/louisbranch/fracturing.space/internal/tools/scenario"
)

// Config holds scenario command configuration.
type Config struct {
	GRPCAddr   string
	Scenario   string
	Assertions bool
	Verbose    bool
	Timeout    time.Duration
}

// ParseConfig parses flags into a Config.
func ParseConfig(fs *flag.FlagSet, args []string) (Config, error) {
	var cfg Config
	fs.StringVar(&cfg.GRPCAddr, "grpc-addr", "localhost:8080", "game server address")
	fs.StringVar(&cfg.Scenario, "scenario", "", "path to scenario lua file")
	fs.BoolVar(&cfg.Assertions, "assert", true, "enable assertions (disable to log expectations)")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "enable verbose logging")
	fs.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "timeout per step")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Run executes the scenario command.
func Run(ctx context.Context, cfg Config, out io.Writer, errOut io.Writer) error {
	if out == nil {
		out = io.Discard
	}
	if errOut == nil {
		errOut = io.Discard
	}
	if cfg.Scenario == "" {
		return errors.New("scenario path is required")
	}

	mode := scenario.AssertionStrict
	if !cfg.Assertions {
		mode = scenario.AssertionLogOnly
	}

	logger := log.New(errOut, "", 0)
	return scenario.RunFile(ctx, scenario.Config{
		GRPCAddr:   cfg.GRPCAddr,
		Timeout:    cfg.Timeout,
		Assertions: mode,
		Verbose:    cfg.Verbose,
		Logger:     logger,
	}, cfg.Scenario)
}
