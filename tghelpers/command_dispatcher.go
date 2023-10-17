package tghelpers

import (
	"context"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/alufers/paczkobot/httphelpers"
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
		go d.processIncomingUpdate(ctx, u)
	}
	return nil
}

func (d *CommandDispatcher) processIncomingUpdate(ctx context.Context, u tgbotapi.Update) {
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

	ctx = context.WithValue(ctx, UpdateContextKey, u)
	ctx = context.WithValue(ctx, ArgsContextKey, args)

	// TODO: move this to a nice hook
	if strings.HasSuffix(args.CommandName, "@har") {
		ctx = httphelpers.WithHarLoggerStorage(ctx)
	}

	for _, hook := range d.UpdateHooks {
		ctx = hook.OnUpdate(ctx)
	}

	shouldProcessCommands := ctx.Value(StopProcessingCommandsCtxKey) == nil
	var err error
	if shouldProcessCommands {
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
	}

	storage := httphelpers.GetHarLoggerStorage(ctx)

	if storage != nil {

		jsonData, err := storage.GetJSONData()
		if err != nil {
			log.Printf("Error while getting HAR data: %v", err)
		} else {
			sendDoc := tgbotapi.NewDocument(args.ChatID, tgbotapi.FileReader{
				Name:   "har.json",
				Reader: strings.NewReader(string(jsonData)),
			})
			sendDoc.Caption = "HAR data"
			_, err := d.BotAPI.Send(sendDoc)
			if err != nil {
				log.Printf("Error while sending HAR data: %v", err)
			}
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
