package tghelpers

import (
	"context"
)

type stopProcessingCommandsCtxKeyType struct{}

// StopProcessingCommandsCtxKey is a context key that can be used to stop
// processing commands for an update.

// It should be added to the context by an UpdateHook, which
// wishes to stop processing commands.
var StopProcessingCommandsCtxKey = stopProcessingCommandsCtxKeyType{}

func WithStopProcessingCommands(ctx context.Context) context.Context {
	return context.WithValue(ctx, StopProcessingCommandsCtxKey, true)
}

// UpdateHook allows a service to listen for all telegram updates
// before they are processed for commands
type UpdateHook interface {
	// OnUpdate is called for each incoming update.
	// If the implementer returns true the update is regarded
	// as handled by the hook. Further processing is stopped.
	// The update shall be extracted from the context using
	// tghelpers.UpdateFromCtx(ctx)
	OnUpdate(context.Context) context.Context

	OnAfterUpdate(context.Context) context.Context
}
