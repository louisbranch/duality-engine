package domain

import (
	"context"
	"time"
)

type toolInvocationContext struct {
	RunCtx       context.Context
	Cancel       context.CancelFunc
	MCPContext   Context
	InvocationID string
}

func newToolInvocationContext(ctx context.Context, getContext func() Context) (toolInvocationContext, error) {
	return newToolInvocationContextWithTimeout(ctx, getContext, grpcCallTimeout)
}

func newToolInvocationContextWithTimeout(ctx context.Context, getContext func() Context, timeout time.Duration) (toolInvocationContext, error) {
	invocationID, err := NewInvocationID()
	if err != nil {
		return toolInvocationContext{}, err
	}

	runCtx, cancel := context.WithTimeout(ctx, timeout)

	mcpCtx := Context{}
	if getContext != nil {
		mcpCtx = getContext()
	}

	return toolInvocationContext{
		RunCtx:       runCtx,
		Cancel:       cancel,
		MCPContext:   mcpCtx,
		InvocationID: invocationID,
	}, nil
}
