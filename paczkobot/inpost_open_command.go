package paczkobot

import (
	"context"
	"fmt"
	"log"

	"github.com/alufers/paczkobot/inpostextra"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InpostOpenCommand struct {
	App *BotApp
}

func (s *InpostOpenCommand) Aliases() []string {
	return []string{"/inpostopen"}
}

func (s *InpostOpenCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{
		{
			Name:        "trackingNumber",
			Description: "Tracking number of the package",
			Question:    "Please enter the tracking number:",
			Variadic:    false,
		},
	}
}

func (f *InpostOpenCommand) Help() string {
	return "Remotely opens an Inpost locker."
}

func (f *InpostOpenCommand) Category() string {
	return "Inpost"
}

func (f *InpostOpenCommand) Execute(ctx context.Context, args *CommandArguments) error {

	suggestions := map[string]string{}
	followedPackages := []FollowedPackageTelegramUser{}

	if err := f.App.DB.Where("telegram_user_id = ?", args.FromUserID).
		Preload("FollowedPackage").
		Preload("FollowedPackage.FollowedPackageProviders").
		Find(&followedPackages).Error; err != nil {
		return fmt.Errorf("failed to query DB: %w", err)
	}

	for _, p := range followedPackages {
		var inpostProvider *FollowedPackageProvider
		for _, prov := range p.FollowedPackage.FollowedPackageProviders {
			if prov.ProviderName == "inpost" {
				inpostProvider = prov
				break
			}
		}
		if inpostProvider == nil || (inpostProvider.LastStatusValue != "ready_to_pickup" && inpostProvider.LastStatusValue != "stack_in_box_machine") {
			continue
		}
		suggestions[p.FollowedPackage.TrackingNumber] = "📦 " + p.CustomName + " (" + p.FollowedPackage.TrackingNumber + ")"
		if p.FollowedPackage.FromName != "" {
			suggestions[p.FollowedPackage.TrackingNumber] += " from " + p.FollowedPackage.FromName
		}
	}

	trackingNumber, err := args.GetOrAskForArgument("trackingNumber", suggestions)
	if err != nil {
		return err
	}

	creds := []*inpostextra.InpostCredentials{}
	if err := f.App.DB.Where("telegram_user_id = ?", args.FromUserID).Find(&creds).Error; err != nil {
		return fmt.Errorf("failed to get inpost credentials: %v", err)
	}
	if len(creds) == 0 {
		return fmt.Errorf("no Inpost credentials found. Please add them using /inpostlogin")
	}

	for _, cred := range creds {
		p, err := f.App.InpostService.GetParcel(f.App.DB, cred, trackingNumber)
		if err != nil {
			log.Printf("failed to get parcel: %v", err)
			continue
		}
		if p.OpenCode != "" && p.PickupPoint != nil {
			err := f.App.AskService.Confirm(args.ChatID, fmt.Sprintf(
				"Do you really want to open the locker at %s (%s %s, %s) for the package with this number: %v?\n\n<b>WARNING:</b> The locker will open immediately, without checking your location. Make sure you are at the correct locker before confirming.",
				p.PickupPoint.Name,
				p.PickupPoint.AddressDetails.Street,
				p.PickupPoint.AddressDetails.BuildingNumber,
				p.PickupPoint.AddressDetails.City,
				trackingNumber,
			),
			)
			if err != nil {
				return err
			}
			err = f.App.InpostService.OpenParcelLocker(f.App.DB, cred, p.ShipmentNumber)
			if err != nil {
				return err
			}
			msg := tgbotapi.NewMessage(args.ChatID, "Locker opened.")
			_, err = f.App.Bot.Send(msg)
			return err
		}
	}

	msg := tgbotapi.NewMessage(args.ChatID, "No account of yours can access the package: "+trackingNumber)
	_, err = f.App.Bot.Send(msg)
	return err
}
