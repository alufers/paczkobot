package tghelpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotAPI is an interface wrapper over *tgbotapi.BotAPI that allows for mocking.
type BotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	GetFile(config tgbotapi.FileConfig) (tgbotapi.File, error)
}
