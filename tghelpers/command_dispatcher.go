package tghelpers

import (
	"context"
	"fmt"
	"html"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CommandDispatcher manages registered commands and
// dispatches incoming messages to them
type CommandDispatcher struct {
	BotAPI      BotAPI
	Commands    []Command
	UpdateHooks []UpdateHook
	AskService  *AskService
}

func NewCommandDispatcher(botAPI BotAPI, AskService *AskService) *CommandDispatcher {
	return &CommandDispatcher{
		BotAPI:      botAPI,
		Commands:    []Command{},
		UpdateHooks: []UpdateHook{},
		AskService:  AskService,
	}
}

func (d *CommandDispatcher) RegisterCommands(commands ...Command) {
	d.Commands = append(d.Commands, commands...)
}

func (d *CommandDispatcher) RegisterUpdateHooks(hooks ...UpdateHook) {
	d.UpdateHooks = append(d.UpdateHooks, hooks...)
}

func (d *CommandDispatcher) RunUpdateLoop() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := d.BotAPI.GetUpdatesChan(u)
	ctx := context.Background()
	for u := range updates {
		// run command processing asynchronously
		go func(u tgbotapi.Update) {
			var cmdText string
			args := &CommandArguments{
				AskService:     d.AskService,
				Update:         &u,
				NamedArguments: map[string]string{},
			}
			if u.Message != nil {
				cmdText = u.Message.Text
				args.ChatID = u.Message.Chat.ID
				args.FromUserID = u.Message.From.ID
			}
			if u.CallbackQuery != nil {
				cmdText = u.CallbackQuery.Data
				args.ChatID = u.CallbackQuery.Message.Chat.ID
				args.FromUserID = u.CallbackQuery.From.ID
			}
			seg := strings.Split(cmdText, " ")
			args.CommandName = seg[0]
			args.Arguments = seg[1:]

			ctx := context.WithValue(ctx, UpdateContextKey, u)
			ctx = context.WithValue(ctx, ArgsContextKey, args)
			for _, hook := range d.UpdateHooks {
				if hook.OnUpdate(ctx) {
					return // hook has handled the message stop processing
				}
			}

			var err error

			for _, cmd := range d.Commands {
				if CommandMatches(cmd, cmdText) {
					args.Command = cmd
					for i, argTpl := range cmd.Arguments() {
						if argTpl.Variadic {
							args.NamedArguments[argTpl.Name] = strings.Join(args.Arguments[i:], " ")
							break
						}
						if i >= len(args.Arguments) {
							break
						}
						args.NamedArguments[argTpl.Name] = args.Arguments[i]
					}

					err = cmd.Execute(ctx)

					break
				}
			}
			if err != nil {
				log.Printf("Error while processing command %v: %v", cmdText, err)
				msg := tgbotapi.NewMessage(args.ChatID, "🚫 Error: <b>"+html.EscapeString(err.Error())+"</b>")
				msg.ParseMode = "HTML"
				if u.Message != nil {
					msg.ReplyToMessageID = u.Message.MessageID
				}
				_, err := d.BotAPI.Send(msg)
				if err != nil {
					log.Printf("An error has occurred while sending an error message: %v", err)
				}
			}
		}(u)
	}
	return nil
}

// RequestSetMyCommands registers the commands in the Telegram bot API.
// so that an autocomplete menu appears when the user types "/"
func (d *CommandDispatcher) RequestSetMyCommands() error {
	myCommands := []tgbotapi.BotCommand{}
	for _, cmd := range d.Commands {
		rawCmd := strings.TrimPrefix(strings.Split(cmd.Aliases()[0], " ")[0], "/")
		myCommands = append(myCommands, tgbotapi.BotCommand{
			Command:     rawCmd,
			Description: cmd.Help(),
		})
	}

	commandsConfig := tgbotapi.NewSetMyCommands(myCommands...)

	if _, err := d.BotAPI.Request(commandsConfig); err != nil {
		return fmt.Errorf("fsailed to set my commands: %v", err)
	}
	return nil
}
