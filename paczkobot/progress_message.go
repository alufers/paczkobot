package paczkobot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ProgressMessage struct {
	Bot *tgbotapi.BotAPI
	Msg *tgbotapi.Message
}

func NewProgressMessage(bot *tgbotapi.BotAPI, chatID int64, initialText string) (*ProgressMessage, error) {
	msg := tgbotapi.NewMessage(chatID, initialText)
	msg.ParseMode = "HTML"

	sentMsg, err := bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &ProgressMessage{
		Bot: bot,
		Msg: &sentMsg,
	}, nil
}

func (p *ProgressMessage) UpdateText(text string) error {
	msg := tgbotapi.NewEditMessageText(p.Msg.Chat.ID, p.Msg.MessageID, text)
	msg.ParseMode = "HTML"

	_, err := p.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProgressMessage) Delete() error {
	msg := tgbotapi.NewDeleteMessage(p.Msg.Chat.ID, p.Msg.MessageID)

	_, err := p.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
