package paczkobot

import (
	"context"
	"fmt"

	"github.com/alufers/paczkobot/inpostextra"
	"github.com/alufers/paczkobot/tghelpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InpostAccountsCommand struct {
	App *BotApp
}

func (s *InpostAccountsCommand) Aliases() []string {
	return []string{"/inpostaccounts"}
}

func (s *InpostAccountsCommand) Arguments() []*tghelpers.CommandDefArgument {
	return []*tghelpers.CommandDefArgument{}
}

func (f *InpostAccountsCommand) Help() string {
	return "Shows inpost accounts you are logged into"
}

func (f *InpostAccountsCommand) Category() string {
	return "Inpost"
}

func (f *InpostAccountsCommand) Execute(ctx context.Context) error {
	args := tghelpers.ArgsFromCtx(ctx)
	creds := []*inpostextra.InpostCredentials{}
	if err := f.App.DB.Where("telegram_user_id = ?", args.FromUserID).Find(&creds).Error; err != nil {
		return fmt.Errorf("failed to get inpost credentials: %v", err)
	}
	msgContent := "Your inpost accounts: \n"
	for _, cred := range creds {
		msgContent += fmt.Sprintf("- <b>%s</b> \n", cred.PhoneNumber)
	}
	if len(creds) == 0 {
		msgContent = "You are not logged into any inpost accounts!"
	}
	msg := tgbotapi.NewMessage(args.ChatID, msgContent)
	msg.ParseMode = "HTML"
	_, err := f.App.Bot.Send(msg)

	return err
}
