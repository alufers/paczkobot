package paczkobot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"html"
	"log"
	"strings"
)

type BotApp struct {
	Bot      *tgbotapi.BotAPI
	Commands []Command
}

func NewBotApp(b *tgbotapi.BotAPI) (a *BotApp) {
	a = &BotApp{
		Bot: b,
	}
	a.Commands = []Command{
		&StartCommand{App: a, ExtraHelp: []Helpable{
			&AvailableProvidersExtraHelp{},
		}},
		&TrackCommand{App: a},
	}
	return
}

func (a *BotApp) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := a.Bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatalf("telegram updates error: %v", err)
	}
	log.Printf("Telegram bot is starting...")

	for u := range updates {
		go func(update tgbotapi.Update) {
			var err error
			log.Printf("msg: %v", update.Message.Text)
			for _, cmd := range a.Commands {
				if CommandMatches(cmd, update.Message.Text) {
					ctx := context.TODO()
					seg := strings.Split(update.Message.Text, " ")
					err = cmd.Execute(ctx, &CommandArguments{
						update:      &update,
						CommandName: seg[0],
						Arguments:   seg[1:],
					})
				}
			}
			log.Print(err)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Error: <b>"+html.EscapeString(err.Error())+"</b>")
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID

				a.Bot.Send(msg)
			}
		}(u)
	}

}
