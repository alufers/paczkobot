package paczkobot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnfollowAllCommand struct {
	App *BotApp
}

func (s *UnfollowAllCommand) Aliases() []string {
	return []string{"/unfollowall"}
}

func (s *UnfollowAllCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{}
}

func (s *UnfollowAllCommand) Help() string {
	return "stops following all packages for changes"
}

func (s *UnfollowAllCommand) Execute(ctx context.Context, args *CommandArguments) error {

	followedPackages := []*FollowedPackageTelegramUser{}

	if err := s.App.DB.Where("telegram_user_id = ?", args.FromUserID).Preload("FollowedPackage").Find(&followedPackages).Error; err != nil {
		return fmt.Errorf("failed to query DB for packages: %v", err)
	}

	err := s.App.AskService.Confirm(args.ChatID, fmt.Sprintf("Are you sure you want to stop following all (%v) packages?", len(followedPackages)))
	if err != nil {
		return fmt.Errorf("Canceled")
	}

	for _, followedPackage := range followedPackages {
		// delete the package from the DB
		if err := s.App.DB.Delete(followedPackage).Error; err != nil {
			return fmt.Errorf("failed to delete package: %v", err)
		}

		// check whether the package is orphaned
		// count FollowedPackageTelegramUser where ID followedPackage.FollowedPackageID
		count := int64(0)
		if err := s.App.DB.Model(&FollowedPackageTelegramUser{}).Where("followed_package_id = ?", followedPackage.FollowedPackageID).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count packages: %v", err)
		}
		if count <= 0 {
			// delete the package from the DB
			if err := s.App.DB.Delete(followedPackage.FollowedPackage).Error; err != nil {
				return fmt.Errorf("failed to delete package: %v", err)
			}
		}
	}

	msg := tgbotapi.NewMessage(args.update.Message.Chat.ID, fmt.Sprintf(`Removed %v followed packages!`, len(followedPackages)))
	msg.ParseMode = "HTML"
	_, err = s.App.Bot.Send(msg)
	return err
}
