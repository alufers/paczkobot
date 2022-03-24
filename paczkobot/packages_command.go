package paczkobot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xeonx/timeago"
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

	packagesText := ""

	followedPackages := []FollowedPackageTelegramUser{}

	if err := s.App.DB.Where("telegram_user_id = ?", args.FromUserID).
		Preload("FollowedPackage").
		Preload("FollowedPackage.FollowedPackageProviders").
		Find(&followedPackages).Error; err != nil {
		return fmt.Errorf("failed to query DB: %w", err)
	}

	for _, p := range followedPackages {
		customName := p.CustomName

		if p.FollowedPackage.FromName != "" {
			if customName != "" {
				customName += " "
			}
			customName = fmt.Sprintf("from %s", p.FollowedPackage.FromName)
		}
		if customName != "" {
			customName = fmt.Sprintf(" <i>(%s)</i>", customName)
		}
		packagesText += fmt.Sprintf("<b>%v</b>%v", p.FollowedPackage.TrackingNumber, customName)
		for i, prov := range p.FollowedPackage.FollowedPackageProviders {
			packagesText += fmt.Sprintf(" %v (<i>%v %v</i>)",
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
