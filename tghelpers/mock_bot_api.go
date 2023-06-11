package tghelpers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MockBotApi struct {
	SendCallback func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	SendCount    int
}

func (m *MockBotApi) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.SendCount++
	return m.SendCallback(c)
}

func (m *MockBotApi) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockBotApi) GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return nil
}

func (m *MockBotApi) GetFile(config tgbotapi.FileConfig) (tgbotapi.File, error) {
	return tgbotapi.File{}, nil
}
