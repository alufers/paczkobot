package tghelpers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// UpdateHook allows a service to listen for all telegram updates
// before they are processed for commands
type UpdateHook interface {
	OnUpdate(update tgbotapi.Update) bool
}
