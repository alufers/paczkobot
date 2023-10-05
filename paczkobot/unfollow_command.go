package paczkobot

import (
	"context"
	"fmt"

	"github.com/alufers/paczkobot/tghelpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnfollowCommand struct {
	App *BotApp
}

func (s *UnfollowCommand) Aliases() []string {
	return []string{"/unfollow"}
}

func (s *UnfollowCommand) Arguments() []*tghelpers.CommandDefArgument {
	return []*tghelpers.CommandDefArgument{
		{
			Name:        "shipmentNumber",
			Description: "shipment number of the package",
			Question:    "Please enter the shipment number to unfollow:",
		},
	}
}

func (s *UnfollowCommand) Help() string {
	return "stops following a package for changes"
}

func (s *UnfollowCommand) Category() string {
	return "Following packages"
}

func (s *UnfollowCommand) Execute(ctx context.Context) error {
	args := tghelpers.ArgsFromCtx(ctx)
	shipmentNumber, err := args.GetOrAskForArgument("shipmentNumber")
	if err != nil {
		return err
	}

	followedPackage := &FollowedPackage{}

	if err := s.App.DB.Where("tracking_number = ?", shipmentNumber).Preload("FollowedPackageTelegramUsers").First(followedPackage).Error; err != nil {
		return fmt.Errorf("no such followed package")
	}

	var currentUser *FollowedPackageTelegramUser
	for _, tgUser := range followedPackage.FollowedPackageTelegramUsers {
		if tgUser.ChatID == args.ChatID {
			currentUser = tgUser
			break
		}
	}

	if currentUser == nil {
		return fmt.Errorf("no such followed package")
	}

	if err := s.App.DB.Delete(currentUser).Error; err != nil {
		return fmt.Errorf("failed to delete followed package")
	}

	if len(followedPackage.FollowedPackageTelegramUsers) <= 1 {
		if err := s.App.DB.Delete(followedPackage).Error; err != nil {
			return fmt.Errorf("failed to delete followed package")
		}
	}

	msg := tgbotapi.NewMessage(args.ChatID, fmt.Sprintf(`Package %v has been unfollowed!`, shipmentNumber))
	msg.ParseMode = "HTML"
	_, err = s.App.Bot.Send(msg)
	return err
}
