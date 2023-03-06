package paczkobot

import (
	"context"
	"fmt"
	"html"
	"sort"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

type StartCommand struct {
	App       *BotApp
	ExtraHelp []Helpable
}

func (s *StartCommand) Aliases() []string {
	return []string{"/start"}
}

func (s *StartCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{}
}

func (s *StartCommand) Help() string {
	return "prints the available commands"
}

func (s *StartCommand) Execute(ctx context.Context, args *CommandArguments) error {
	categoriesHelp := map[string][]Command{}
	for _, cmd := range s.App.Commands {
		if cmdWithCat, ok := cmd.(CommandWithCategory); ok {
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
			if helpable, ok := cmd.(Helpable); ok {
				line += " - " + html.EscapeString(helpable.Help())
			}
			commandHelp += line + "\n"
		}
		commandHelp += "\n"
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
