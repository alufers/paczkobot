package paczkobot

import (
	"context"
	"fmt"

	"github.com/alufers/paczkobot/inpostextra"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InpostLogoutCommand struct {
	App *BotApp
}

func (s *InpostLogoutCommand) Aliases() []string {
	return []string{"/inpostlogout"}
}

func (s *InpostLogoutCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{
		{
			Name:        "phoneNumber",
			Description: "The phone number associated with the Inpost account to logout from",
			Question:    "Please enter the phone number to logout from:",
			Variadic:    false,
		},
	}
}

func (f *InpostLogoutCommand) Help() string {
	return "logs out from the inpost service"
}

func (f *InpostLogoutCommand) Category() string {
	return "Inpost"
}

func (f *InpostLogoutCommand) Execute(ctx context.Context, args *CommandArguments) error {
	creds := []*inpostextra.InpostCredentials{}
	if err := f.App.DB.Where("telegram_user_id = ?", args.FromUserID).Find(&creds).Error; err != nil {
		return fmt.Errorf("failed to get inpost credentials: %v", err)
	}
	suggestions := map[string]string{}
	for _, cred := range creds {
		suggestions[cred.PhoneNumber] = cred.PhoneNumber
	}
	phoneNumber, err := args.GetOrAskForArgument("phoneNumber", suggestions)
	if err != nil {
		return err
	}
	// Delete the credentials from the database
	err = f.App.DB.Where("telegram_user_id = ? AND phone_number = ?", args.FromUserID, phoneNumber).Delete(&inpostextra.InpostCredentials{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete existing credentials: %v", err)
	}

	msg := tgbotapi.NewMessage(args.ChatID, "Successfully logged out from Inpost")
	if _, err := f.App.Bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}
