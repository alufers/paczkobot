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
		&SetNameCommand{App: a},
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
	a.Bot.SetChatDescription(tgbotapi.SetChatDescriptionConfig{})
	go a.TrackingService.RunAutomaticTrackingLoop()
	for u := range updates {
		go func(update tgbotapi.Update) {
			var err error
			var cmdText string

			args := &CommandArguments{
				update: &update,
			}
			if update.Message != nil {
				cmdText = update.Message.Text
				args.ChatID = update.Message.Chat.ID
				args.FromUserID = update.Message.From.ID
			}
			if update.CallbackQuery != nil {
				cmdText = update.CallbackQuery.Data
				args.ChatID = update.CallbackQuery.Message.Chat.ID
				args.FromUserID = update.CallbackQuery.From.ID
			}
			seg := strings.Split(cmdText, " ")
			args.CommandName = seg[0]
			args.Arguments = seg[1:]

			for _, cmd := range a.Commands {

				if CommandMatches(cmd, cmdText) {
					ctx := context.TODO()
					err = cmd.Execute(ctx, args)
				}
			}

			if err != nil {
				log.Printf("Error while processing command %v: %v", cmdText, err)
				msg := tgbotapi.NewMessage(args.ChatID, "🚫 Error: <b>"+html.EscapeString(err.Error())+"</b>")
				msg.ParseMode = "HTML"
				if update.Message != nil {
					msg.ReplyToMessageID = update.Message.MessageID
				}
				a.Bot.Send(msg)
			}
		}(u)
	}

}
