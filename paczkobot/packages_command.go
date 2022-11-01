package paczkobot

import (
	"context"
	"fmt"
	"log"

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


func (s *PackagesCommand) Category() string {
	return "Following packages"
}


func (s *PackagesCommand) Execute(ctx context.Context, args *CommandArguments) error {

	if err := s.App.ArchiveService.FetchAndArchivePackagesForUser(args.FromUserID); err != nil {
		log.Printf("failed to fetch and archive packages for user %v: %v", args.FromUserID, err)
	}

	followedPackages := []FollowedPackageTelegramUser{}

	if err := s.App.DB.Where("telegram_user_id = ? AND archived = ?", args.FromUserID, false).
		Preload("FollowedPackage").
		Preload("FollowedPackage.FollowedPackageProviders").
		Find(&followedPackages).Error; err != nil {
		return fmt.Errorf("failed to query DB: %w", err)
	}

	msg := tgbotapi.NewMessage(args.ChatID, fmt.Sprintf(`
Your followed packages:

%v

Use /archived to see your archived packages.
`, s.App.PackagePrinterService.PrintPackages(followedPackages)))
	msg.ParseMode = "HTML"
	_, err := s.App.Bot.Send(msg)
	return err
}
