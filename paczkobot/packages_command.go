package paczkobot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PackagesCommand struct {
	App *BotApp
}

func (s *PackagesCommand) Aliases() []string {
	return []string{"/packages"}
}

func (s *PackagesCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{}
}

func (s *PackagesCommand) Help() string {
	return "prints your followed packages"
}

func (s *PackagesCommand) Execute(ctx context.Context, args *CommandArguments) error {

	followedPackages := []FollowedPackageTelegramUser{}

	if err := s.App.DB.Where("telegram_user_id = ?", args.FromUserID).
		Preload("FollowedPackage").
		Preload("FollowedPackage.FollowedPackageProviders").
		Find(&followedPackages).Error; err != nil {
		return fmt.Errorf("failed to query DB: %w", err)
	}

	msg := tgbotapi.NewMessage(args.ChatID, fmt.Sprintf(`
Your followed packages:

%v

`, s.App.PackagePrinterService.PrintPackages(followedPackages)))
	msg.ParseMode = "HTML"
	_, err := s.App.Bot.Send(msg)
	return err
}
