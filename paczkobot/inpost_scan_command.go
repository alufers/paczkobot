package paczkobot

import (
	"context"
	"fmt"

	"github.com/alufers/paczkobot/inpostextra"
	"github.com/alufers/paczkobot/tghelpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InpostScanCommand struct {
	App *BotApp
}

func (s *InpostScanCommand) Aliases() []string {
	return []string{"/inpostscan"}
}

func (s *InpostScanCommand) Arguments() []*tghelpers.CommandDefArgument {
	return []*tghelpers.CommandDefArgument{}
}

func (f *InpostScanCommand) Help() string {
	return "Scans all your inpost accounts and follows any new packages"
}

func (f *InpostScanCommand) Category() string {
	return "Inpost"
}

func (f *InpostScanCommand) Execute(ctx context.Context) error {
	args := tghelpers.ArgsFromCtx(ctx)
	creds := []*inpostextra.InpostCredentials{}
	if err := f.App.DB.Where("telegram_user_id = ?", args.FromUserID).Find(&creds).Error; err != nil {
		return fmt.Errorf("failed to get inpost credentials: %v", err)
	}
	errors := []error{}
	successCount := 0
	for _, cred := range creds {
		err := f.App.InpostScannerService.ScanUserPackages(
			ctx,
			cred,
		)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to scan inpost account %v: %v", cred.PhoneNumber, err))
		} else {
			successCount++
		}
	}
	msgTxt := fmt.Sprintf(`Successfully scanned %v out of %v inpost accounts!`, successCount, len(creds))
	if len(errors) > 0 {
		msgTxt += " Errors: \n"
		for _, err := range errors {
			msgTxt += fmt.Sprintf("- %v \n", err)
		}
	}
	msg := tgbotapi.NewMessage(args.Update.Message.Chat.ID, msgTxt)
	msg.ParseMode = "HTML"
	_, err := f.App.Bot.Send(msg)

	return err
}
