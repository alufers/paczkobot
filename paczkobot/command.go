package paczkobot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type Command interface {
	Usage() string
	Execute(ctx context.Context, args *CommandArguments) error
}

type Helpable interface {
	Help() string
}

type CommandArguments struct {
	update      *tgbotapi.Update
	CommandName string
	Arguments   []string
}

func CommandMatches(cmd Command, userInput string) bool {
	usage := cmd.Usage()
	return strings.Split(userInput, " ")[0] == strings.Split(usage, " ")[0]
}
