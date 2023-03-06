package paczkobot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InpostLoginCommand struct {
	App *BotApp
}

func (s *InpostLoginCommand) Aliases() []string {
	return []string{"/inpostlogin"}
}

func (s *InpostLoginCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{}
}

func (f *InpostLoginCommand) Help() string {
	return "logs in to the inpost service using a sms code"
}

func (f *InpostLoginCommand) Category() string {
	return "Inpost"
}

func (f *InpostLoginCommand) Execute(ctx context.Context, args *CommandArguments) error {
	phoneNumber, err := f.App.AskService.AskForArgument(args.ChatID, "Enter your phone number associated with your inpost account:")
	if err != nil {
		return err
	}
	err = f.App.InpostService.SendSMSCode(phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to send sms code: %v", err)
	}
	code, err := f.App.AskService.AskForArgument(args.ChatID, fmt.Sprintf("Enter the sms code sent to %s:", phoneNumber))
	if err != nil {
		return err
	}

	creds, err := f.App.InpostService.ConfirmSMSCode(phoneNumber, code)
	if err != nil {
		return fmt.Errorf("failed to confirm sms code: %v", err)
	}

	creds.TelegramUserID = args.FromUserID
	err = f.App.DB.Where("telegram_user_id = ? AND phone_number = ?", args.FromUserID, phoneNumber).FirstOrCreate(&creds).Error
	if err != nil {
		return fmt.Errorf("failed to delete existing credentials: %v", err)
	}
	if err := f.App.DB.Save(creds).Error; err != nil {
		return fmt.Errorf("failed to save credentials: %v", err)
	}
	err = f.App.InpostScannerService.ScanUserPackages(
		creds,
	)
	if err != nil {
		return fmt.Errorf("failed to scan user inpost packages: %v", err)
	}
	msg := tgbotapi.NewMessage(args.ChatID, fmt.Sprintf("Successfully logged in to Inpost with phone number %s", phoneNumber))
	if _, err := f.App.Bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}
