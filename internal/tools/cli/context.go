// Package cli provides shared helpers for tool CLIs.
package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WithSignalContext returns a context canceled on SIGINT or SIGTERM.
func WithSignalContext(parent context.Context) (context.Context, func()) {
	if parent == nil {
		parent = context.Background()
	}
	return signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
}

// WithSignalTimeout returns a context canceled on SIGINT, SIGTERM, or timeout.
func WithSignalTimeout(parent context.Context, timeout time.Duration) (context.Context, func()) {
	ctx, stop := WithSignalContext(parent)
	if timeout <= 0 {
		return ctx, stop
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	return ctx, func() {
		cancel()
		stop()
	}
}
