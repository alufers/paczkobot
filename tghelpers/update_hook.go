package tghelpers

import (
	"context"
)

// UpdateHook allows a service to listen for all telegram updates
// before they are processed for commands
type UpdateHook interface {
	// OnUpdate is called for each incoming update.
	// If the implementer returns true the update is regarded
	// as handled by the hook. Further processing is stopped.
	// The update shall be extracted from the context using
	// tghelpers.UpdateFromCtx(ctx)
	OnUpdate(context.Context) bool
}
