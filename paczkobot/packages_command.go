package paczkobot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xeonx/timeago"
)

type PackagesCommand struct {
	App *BotApp
}

func (s *PackagesCommand) Usage() string {
	return "/packages"
}

func (s *PackagesCommand) Help() string {
	return "prints your followed commands"
}

func (s *PackagesCommand) Execute(ctx context.Context, args *CommandArguments) error {

	packagesText := ""

	followedPackages := []FollowedPackageTelegramUser{}

	if err := s.App.DB.Where("telegram_user_id = ?", args.FromUserID).
		Preload("FollowedPackage").
		Preload("FollowedPackage.FollowedPackageProviders").
		Find(&followedPackages).Error; err != nil {
		return fmt.Errorf("failed to query DB: %w", err)
	}

	for _, p := range followedPackages {
		customName := ""
		if p.CustomName != "" {
			customName = fmt.Sprintf(" (<i>%s</i>)", p.CustomName)
		}
		packagesText += fmt.Sprintf("<b>%v</b>%v", p.FollowedPackage.TrackingNumber, customName)
		for i, prov := range p.FollowedPackage.FollowedPackageProviders {
			packagesText += fmt.Sprintf("%v (<i>%v %v</i>)",
				prov.ProviderName,
				prov.LastStatusValue,
				timeago.English.Format(prov.LastStatusDate))
			if i != len(p.FollowedPackage.FollowedPackageProviders)-1 {
				packagesText += ", "
			}
		}

		packagesText += "\n"
	}

	msg := tgbotapi.NewMessage(args.ChatID, fmt.Sprintf(`
Your followed packages:

%v

`, packagesText))
	msg.ParseMode = "HTML"
	_, err := s.App.Bot.Send(msg)
	return err
}
