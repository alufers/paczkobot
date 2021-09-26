package paczkobot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UnfollowCommand struct {
	App *BotApp
}

func (s *UnfollowCommand) Usage() string {
	return "/unfollow <shipmentNumber>"
}

func (s *UnfollowCommand) Help() string {
	return "stops following a package for changes"
}

func (s *UnfollowCommand) Execute(ctx context.Context, args *CommandArguments) error {

	var segments = strings.Split(args.update.Message.Text, " ")

	if len(segments) < 2 {
		return fmt.Errorf("usage: /unfollow &lt;shipmentNumber&gt;")
	}
	shipmentNumber := segments[1]

	followedPackage := &FollowedPackage{}

	if err := s.App.DB.Where("tracking_number = ?", shipmentNumber).Preload("FollowedPackageTelegramUsers").First(followedPackage).Error; err != nil {
		return fmt.Errorf("no such followed package")
	}

	var currentUser *FollowedPackageTelegramUser
	for _, tgUser := range followedPackage.FollowedPackageTelegramUsers {
		if tgUser.TelegramUserID == args.update.Message.From.ID {
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

	msg := tgbotapi.NewMessage(args.update.Message.Chat.ID, fmt.Sprintf(`Package %v has been unfollowed!`, shipmentNumber))
	msg.ParseMode = "HTML"
	_, err := s.App.Bot.Send(msg)
	return err
}
