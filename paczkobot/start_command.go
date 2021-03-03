package paczkobot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"html"
)

type StartCommand struct {
	App       *BotApp
	ExtraHelp []Helpable
}

func (s *StartCommand) Usage() string {
	return "/start"
}

func (s *StartCommand) Help() string {
	return "prints the available commands"
}

func (s *StartCommand) Execute(ctx context.Context, args *CommandArguments) error {

	commandHelp := ""
	extraHelp := ""

	for _, cmd := range s.App.Commands {
		line := html.EscapeString(cmd.Usage())
		if helpable, ok := cmd.(Helpable); ok {
			line += " - " + html.EscapeString(helpable.Help())
		}
		commandHelp += line + "\n"
	}

	for _, e := range s.ExtraHelp {
		extraHelp += e.Help() + "\n"
	}

	msg := tgbotapi.NewMessage(args.update.Message.Chat.ID, fmt.Sprintf(`
<b>Welcome to @%v!</b>

Available commands:
%v

%v
`, viper.GetString("telegram.username"), commandHelp, extraHelp))
	msg.ParseMode = "HTML"
	_, err := s.App.Bot.Send(msg)
	return err
}
