package tghelpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotAPI is an interface wrapper over *tgbotapi.BotAPI that allows for mocking.
type BotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}
