package paczkobot

import (
	"context"
	"fmt"
	"html"
	"sort"

	"github.com/alufers/paczkobot/tghelpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

type StartCommand struct {
	App       *BotApp
	ExtraHelp []tghelpers.Helpable
}

func (s *StartCommand) Aliases() []string {
	return []string{"/start"}
}

func (s *StartCommand) Arguments() []*tghelpers.CommandDefArgument {
	return []*tghelpers.CommandDefArgument{}
}

func (s *StartCommand) Help() string {
	return "prints the available commands"
}

func (s *StartCommand) Execute(ctx context.Context) error {
	args := tghelpers.ArgsFromCtx(ctx)
	categoriesHelp := map[string][]tghelpers.Command{}
	for _, cmd := range s.App.CommandDispatcher.Commands {
		if cmdWithCat, ok := cmd.(tghelpers.CommandWithCategory); ok {
			categoriesHelp[cmdWithCat.Category()] = append(categoriesHelp[cmdWithCat.Category()], cmd)
		} else {
			categoriesHelp["Misc"] = append(categoriesHelp["Misc"], cmd)
		}
	}

	categoryKeys := []string{}
	for k := range categoriesHelp {
		categoryKeys = append(categoryKeys, k)
	}
	sort.Strings(categoryKeys)

	commandHelp := ""
	extraHelp := ""

	for _, category := range categoryKeys {
		cmds := categoriesHelp[category]
		commandHelp += fmt.Sprintf("<b>%s</b>:\n", category)
		for _, cmd := range cmds {

			line := html.EscapeString(cmd.Aliases()[0])
			for _, arg := range cmd.Arguments() {
				line += fmt.Sprintf(" [%s]", arg.Name)
			}
			if helpable, ok := cmd.(tghelpers.Helpable); ok {
				line += " - " + html.EscapeString(helpable.Help())
			}
			commandHelp += line + "\n"
		}
		commandHelp += "\n"
	}
	for _, e := range s.ExtraHelp {
		extraHelp += e.Help() + "\n"
	}

	msg := tgbotapi.NewMessage(args.Update.Message.Chat.ID, fmt.Sprintf(`
<b>Welcome to @%v!</b>

Available commands:
%v
%v
`, viper.GetString("telegram.username"), commandHelp, extraHelp))
	msg.ParseMode = "HTML"
	_, err := s.App.Bot.Send(msg)
	return err
}
