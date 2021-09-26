package paczkobot

import (
	"context"
	"html"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

type BotApp struct {
	Bot                  *tgbotapi.BotAPI
	DB                   *gorm.DB
	Commands             []Command
	NotificationsService *NotificationsService
	TrackingService      *TrackingService
}

func NewBotApp(b *tgbotapi.BotAPI, DB *gorm.DB) (a *BotApp) {
	a = &BotApp{
		Bot: b,
		DB:  DB,
	}
	a.Commands = []Command{
		&StartCommand{App: a, ExtraHelp: []Helpable{
			&AvailableProvidersExtraHelp{},
		}},
		&TrackCommand{App: a},
		&FollowCommand{App: a},
		&PackagesCommand{App: a},
		&UnfollowCommand{App: a},
	}
	a.NotificationsService = NewNotificationsService(a)
	a.TrackingService = NewTrackingService(a)
	return
}

func (a *BotApp) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	log.Printf("Flushing enqueued notifications...")
	if err := a.NotificationsService.FlushEnqueuedNotifications(); err != nil {
		log.Fatalf("Failed to flush enqueued notifications: %v", err)
	}
	log.Printf("Done flushing enqueued notifications!")

	updates, err := a.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("telegram updates error: %v", err)
	}
	log.Printf("Telegram bot is starting...")

	myCommands := []tgbotapi.BotCommand{}
	for _, cmd := range a.Commands {
		rawCmd := strings.TrimPrefix(strings.Split(cmd.Usage(), " ")[0], "/")
		myCommands = append(myCommands, tgbotapi.BotCommand{
			Command:     rawCmd,
			Description: cmd.Help(),
		})
	}
	if err := a.Bot.SetMyCommands(myCommands); err != nil {
		log.Fatalf("Failed to set my commands: %v", err)
	}

	for u := range updates {
		go func(update tgbotapi.Update) {
			var err error

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

			if err != nil {
				log.Printf("Error while processing command %v: %v", update.Message.Text, err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Error: <b>"+html.EscapeString(err.Error())+"</b>")
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID

				a.Bot.Send(msg)
			}
		}(u)
	}

}
