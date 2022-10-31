package paczkobot

import (
	"context"
	"fmt"
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
	if err := f.App.DB.Save(creds).Error; err != nil {
		return fmt.Errorf("failed to save credentials: %v", err)
	}
	err = f.App.InpostScannerService.ScanUserPackages(
		creds,
	)
	if err != nil {
		return fmt.Errorf("failed to scan user inpost packages: %v", err)
	}
	return nil
}
