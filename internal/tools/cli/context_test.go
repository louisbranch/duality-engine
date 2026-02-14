package cli

import (
	"context"
	"testing"
	"time"
)

func TestWithSignalContextNilParent(t *testing.T) {
	ctx, stop := WithSignalContext(nil)
	t.Cleanup(stop)
	if ctx == nil {
		t.Fatal("expected context")
	}
}

func TestWithSignalContextStopCancels(t *testing.T) {
	ctx, stop := WithSignalContext(context.Background())
	stop()
	waitForDone(t, ctx, 50*time.Millisecond)
}

func TestWithSignalTimeoutCancels(t *testing.T) {
	ctx, stop := WithSignalTimeout(context.Background(), 5*time.Millisecond)
	t.Cleanup(stop)
	waitForDone(t, ctx, 200*time.Millisecond)
}

func waitForDone(t *testing.T, ctx context.Context, limit time.Duration) {
	t.Helper()
	select {
	case <-ctx.Done():
		return
	case <-time.After(limit):
		t.Fatal("expected context to be done")
	}
}
