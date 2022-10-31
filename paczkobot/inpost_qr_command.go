package paczkobot

import (
	"context"
	"fmt"
	"log"

	"github.com/alufers/paczkobot/inpostextra"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	qrcode "github.com/skip2/go-qrcode"
)

type InpostQrCommand struct {
	App *BotApp
}

func (s *InpostQrCommand) Aliases() []string {
	return []string{"/inpostqr"}
}

func (s *InpostQrCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{
		{
			Name:        "trackingNumber",
			Description: "Tracking number of the package",
			Question:    "Please enter the tracking number:",
			Variadic:    false,
		},
	}
}

func (f *InpostQrCommand) Help() string {
	return "Shows a QR code for a package"
}

func (f *InpostQrCommand) Category() string {
	return "Inpost"
}

func (f *InpostQrCommand) Execute(ctx context.Context, args *CommandArguments) error {

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
		if inpostProvider == nil || inpostProvider.LastStatusValue == "delivered" {
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
		if p.OpenCode != "" {
			var png []byte
			png, err := qrcode.Encode(p.QrCode, qrcode.High, 256)

			photo := tgbotapi.NewPhoto(args.update.Message.Chat.ID, tgbotapi.FileBytes{Name: "qr.png", Bytes: png})
			photo.Caption = fmt.Sprintf(
				"Phone number: <b>%v</b>\nOpen code: <b>%v</b>",
				cred.PhoneNumber,
				p.OpenCode,
			)

			photo.ParseMode = "HTML"
			_, err = f.App.Bot.Send(photo)
			return err
		}
	}

	msg := tgbotapi.NewMessage(args.ChatID, "No account of yours can access the package: "+trackingNumber)
	f.App.Bot.Send(msg)
	return err
}
