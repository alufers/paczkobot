package paczkobot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ArchivedCommand struct {
	App *BotApp
}

func (s *ArchivedCommand) Aliases() []string {
	return []string{"/archived"}
}

func (s *ArchivedCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{}
}

func (s *ArchivedCommand) Help() string {
	return "prints your archived packages"
}

func (s *ArchivedCommand) Execute(ctx context.Context, args *CommandArguments) error {
	if err := s.App.ArchiveService.FetchAndArchivePackagesForUser(args.FromUserID); err != nil {
		log.Printf("failed to fetch and archive packages for user %v: %v", args.FromUserID, err)
	}

	followedPackages := []FollowedPackageTelegramUser{}

	if err := s.App.DB.Where("telegram_user_id = ? AND archived = ?", args.FromUserID, true).
		Preload("FollowedPackage").
		Preload("FollowedPackage.FollowedPackageProviders").
		Find(&followedPackages).Error; err != nil {
		return fmt.Errorf("failed to query DB: %w", err)
	}

	msg := tgbotapi.NewMessage(args.ChatID, fmt.Sprintf(`
Your archived packages:

%v

`, s.App.PackagePrinterService.PrintPackages(followedPackages)))
	msg.ParseMode = "HTML"
	_, err := s.App.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to list %d packages: %w", len(followedPackages), err)
	}
	return nil
}
