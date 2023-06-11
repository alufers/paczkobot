package tghelpers

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UpdateHook allows a service to listen for all telegram updates
// before they are processed for commands
type UpdateHook interface {
	// OnUpdate is called for each incoming update.
	// If the implementer returns true the update is regarded
	// as handled by the hook. Further processing is stopped.
	OnUpdate(context.Context, tgbotapi.Update) bool
}
