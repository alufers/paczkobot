package paczkobot

import (
	"context"
	"fmt"

	"github.com/alufers/paczkobot/inpostextra"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InpostScanCommand struct {
	App *BotApp
}

func (s *InpostScanCommand) Aliases() []string {
	return []string{"/inpostscan"}
}

func (s *InpostScanCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{}
}

func (f *InpostScanCommand) Help() string {
	return "Scans all your inpost accounts and follows any new packages"
}

func (f *InpostScanCommand) Execute(ctx context.Context, args *CommandArguments) error {

	creds := []*inpostextra.InpostCredentials{}
	if err := f.App.DB.Where("telegram_user_id = ?", args.FromUserID).Find(&creds).Error; err != nil {
		return fmt.Errorf("failed to get inpost credentials: %v", err)
	}
	for _, cred := range creds {
		err := f.App.InpostScannerService.ScanUserPackages(
			cred,
		)
		if err != nil {
			return fmt.Errorf("failed to scan user inpost packages: %v", err)
		}
	}
	msg := tgbotapi.NewMessage(args.update.Message.Chat.ID, fmt.Sprintf(`Scanning %v inpost accounts!`, len(creds)))
	msg.ParseMode = "HTML"
	_, err := f.App.Bot.Send(msg)

	return err
}
