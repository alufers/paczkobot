package paczkobot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SetNameCommand struct {
	App *BotApp
}

func (s *SetNameCommand) Aliases() []string {
	return []string{"/setname"}
}

func (s *SetNameCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{
		{
			Name:        "shipmentNumber",
			Description: "shipment number of the package",
			Question:    "Please enter the shipment number to track:",
		},
		{
			Name:        "name",
			Description: "your custom name for the package",
			Question:    "Please enter your custom name for the package:",
		},
	}
}

func (s *SetNameCommand) Help() string {
	return "sets a name for a package"
}

func (s *SetNameCommand) Execute(ctx context.Context, args *CommandArguments) error {

	shipmentNumber, err := args.GetOrAskForArgument("shipmentNumber")
	if err != nil {
		return err
	}

	followedPackage := &FollowedPackage{}

	if err := s.App.DB.Where("tracking_number = ?", shipmentNumber).Preload("FollowedPackageTelegramUsers").First(followedPackage).Error; err != nil {
		return fmt.Errorf("no such followed package")
	}

	for _, tgUser := range followedPackage.FollowedPackageTelegramUsers {
		if tgUser.TelegramUserID == args.FromUserID {
			customName, err := args.GetOrAskForArgument("name")
			if err != nil {
				return err
			}
			tgUser.CustomName = customName
			if err := s.App.DB.Save(tgUser).Error; err != nil {
				return err
			}
			msg := tgbotapi.NewMessage(args.ChatID,
				fmt.Sprintf(`Package %v has been renamed to <b>%v</b>!`, shipmentNumber, tgUser.CustomName),
			)
			msg.ParseMode = "HTML"
			_, err = s.App.Bot.Send(msg)
			return err

		}
	}

	return fmt.Errorf("no such followed package")
}
