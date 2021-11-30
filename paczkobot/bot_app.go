package paczkobot

import (
	"context"
	"html"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type BotApp struct {
	Bot                  *tgbotapi.BotAPI
	DB                   *gorm.DB
	Commands             []Command
	NotificationsService *NotificationsService
	TrackingService      *TrackingService
	AskService           *AskService
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
		&UnfollowAllCommand{App: a},
	}
	a.NotificationsService = NewNotificationsService(a)
	a.TrackingService = NewTrackingService(a)
	a.AskService = NewAskService(a)
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

	updates := a.Bot.GetUpdatesChan(u)

	log.Printf("Telegram bot is starting...")

	myCommands := []tgbotapi.BotCommand{}
	for _, cmd := range a.Commands {
		rawCmd := strings.TrimPrefix(strings.Split(cmd.Aliases()[0], " ")[0], "/")
		myCommands = append(myCommands, tgbotapi.BotCommand{
			Command:     rawCmd,
			Description: cmd.Help(),
		})
	}

	commandsConfig := tgbotapi.NewSetMyCommands(myCommands...)

	if _, err := a.Bot.Request(commandsConfig); err != nil {
		log.Fatalf("Failed to set my commands: %v", err)
	}

	go a.TrackingService.RunAutomaticTrackingLoop()
	for u := range updates {
		go func(update tgbotapi.Update) {
			if a.AskService.ProcessIncomingMessage(update) {
				return
			}
			var err error
			var cmdText string

			args := &CommandArguments{
				BotApp:         a,
				update:         &update,
				namedArguments: map[string]string{},
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
					args.Command = cmd
					for i, argTpl := range cmd.Arguments() {
						if argTpl.Variadic {
							args.namedArguments[argTpl.Name] = strings.Join(args.Arguments[i:], " ")
							break
						}
						if i >= len(args.Arguments) {
							break
						}
						args.namedArguments[argTpl.Name] = args.Arguments[i]
					}
					ctx := context.TODO()
					err = cmd.Execute(ctx, args)
					break
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
