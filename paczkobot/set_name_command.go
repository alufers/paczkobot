package paczkobot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SetNameCommand struct {
	App *BotApp
}

func (s *SetNameCommand) Usage() string {
	return "/setname <shipmentNumber> <name>"
}

func (s *SetNameCommand) Help() string {
	return "sets a name for a package"
}

func (s *SetNameCommand) Execute(ctx context.Context, args *CommandArguments) error {

	if len(args.Arguments) < 2 {
		return fmt.Errorf("usage: /setname <shipmentNumber> <name>")
	}
	shipmentNumber := args.Arguments[0]

	followedPackage := &FollowedPackage{}

	if err := s.App.DB.Where("tracking_number = ?", shipmentNumber).Preload("FollowedPackageTelegramUsers").First(followedPackage).Error; err != nil {
		return fmt.Errorf("no such followed package")
	}

	for _, tgUser := range followedPackage.FollowedPackageTelegramUsers {
		if tgUser.TelegramUserID == args.update.Message.From.ID {
			tgUser.CustomName = strings.Join(args.Arguments[1:], " ")
			if err := s.App.DB.Save(tgUser).Error; err != nil {
				return err
			}
			msg := tgbotapi.NewMessage(args.update.Message.Chat.ID,
				fmt.Sprintf(`Package %v has been renamed to <b>%v</b>!`, shipmentNumber, tgUser.CustomName),
			)
			msg.ParseMode = "HTML"
			_, err := s.App.Bot.Send(msg)
			return err

		}
	}

	return fmt.Errorf("no such followed package")
}
