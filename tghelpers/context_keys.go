package tghelpers

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type contextKeyType string

const (
	UpdateContextKey contextKeyType = "update"
	ArgsContextKey   contextKeyType = "args"
)

func UpdateFromCtx(ctx context.Context) tgbotapi.Update {
	return ctx.Value(UpdateContextKey).(tgbotapi.Update)
}

func ArgsFromCtx(ctx context.Context) *CommandArguments {
	return ctx.Value(ArgsContextKey).(*CommandArguments)
}
