package tghelpers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MockBotApi struct {
	SendCallback func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	SendCount    int
}

func (m *MockBotApi) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.SendCount++
	return m.SendCallback(c)
}
