package paczkobot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/providers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FollowCommand struct {
	App *BotApp
}

func (s *FollowCommand) Aliases() []string {
	return []string{"/follow"}
}

func (s *FollowCommand) Arguments() []*CommandDefArgument {
	return []*CommandDefArgument{
		&CommandDefArgument{
			Name:        "shipmentNumber",
			Description: "shipment number of the package",
			Question:    "Please enter the shipment number:",
		},
	}
}

func (f *FollowCommand) Help() string {
	return "follows a package and sends you an update every time its status changes"
}

func (f *FollowCommand) Execute(ctx context.Context, args *CommandArguments) error {
	msg := tgbotapi.NewMessage(args.ChatID, "⌛ loading...")
	msg.ParseMode = "HTML"
	loadingRes, err := f.App.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send loading message: %w", err)
	}
	defer func() {
		f.App.Bot.Send(tgbotapi.NewDeleteMessage(args.ChatID, loadingRes.MessageID))
	}()

	shipmentNumber, err := args.GetOrAskForArgument("shipmentNumber")
	if err != nil {
		return err
	}
	log.Printf("following shipmentNumber = %v", shipmentNumber)
	providersToCheck := []providers.Provider{}
	for _, provider := range providers.AllProviders {
		if provider.MatchesNumber(shipmentNumber) {
			providersToCheck = append(providersToCheck, provider)
		}
	}

	if len(providersToCheck) == 0 {
		return fmt.Errorf("no tracking providers support this tracking number")
	}

	providersToFollow := []*FollowedPackageProvider{}

	for _, p := range providersToCheck {

		d, err := providers.InvokeProvider(context.Background(), p, shipmentNumber)
		if errors.Is(err, commonerrors.NotFoundError) {
			continue
		}
		if err != nil {
			continue
		}
		lastStep := &commondata.TrackingStep{}
		if len(d.TrackingSteps) > 0 {
			lastStep = d.TrackingSteps[len(d.TrackingSteps)-1]
		}
		providersToFollow = append(providersToFollow, &FollowedPackageProvider{
			ProviderName:       p.GetName(),
			LastStatusValue:    lastStep.Message,
			LastStatusDate:     lastStep.Datetime,
			LastStatusLocation: lastStep.Location,
		})
	}

	if len(providersToFollow) <= 0 {
		return fmt.Errorf("the package was not found in any of the providers, please double check your package number")
	}

	followedPackage := &FollowedPackage{
		TrackingNumber:           shipmentNumber,
		LastAutomaticCheck:       time.Now(),
		LastChange:               time.Now(),
		FollowedPackageProviders: providersToFollow,
	}

	if err := f.App.DB.Unscoped().Where("tracking_number = ?", shipmentNumber).FirstOrCreate(followedPackage).Error; err != nil {
		return fmt.Errorf("failed to create FollowedPackage: %v", err)
	}

	if followedPackage.DeletedAt.Valid {
		followedPackage.DeletedAt.Valid = false
		if err := f.App.DB.Save(followedPackage).Error; err != nil {
			return fmt.Errorf("failed to restore FollowedPackage: %v", err)
		}
	}

	followedPackageTelegramUser := &FollowedPackageTelegramUser{
		FollowedPackageID: followedPackage.ID,
		TelegramUserID:    args.FromUserID,
		ChatID:            args.ChatID,
	}

	if err := f.App.DB.Where("followed_package_id = ? AND telegram_user_id = ?",
		followedPackage.ID,
		followedPackageTelegramUser.TelegramUserID,
	).FirstOrCreate(followedPackageTelegramUser).Error; err != nil {
		return fmt.Errorf("failed to create FollowedPackageTelegramUser: %v", err)
	}

	for _, p := range providersToFollow {
		p.FollowedPackageID = followedPackage.ID
		if err := f.App.DB.Where("followed_package_id = ? AND provider_name = ?",
			followedPackage.ID,
			p.ProviderName,
		).FirstOrCreate(p).Error; err != nil {
			return fmt.Errorf("failed to create FollowedPackageTelegramUser: %v", err)
		}
	}
	msg2 := tgbotapi.NewMessage(args.ChatID, fmt.Sprintf(`Package %v has been added to your followed packages!`, shipmentNumber))
	msg2.ParseMode = "HTML"
	msg2.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🖊️ Set a name for this package", "/setname "+shipmentNumber),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📦 See followed packages", "/packages"),
		),
	)
	_, err = f.App.Bot.Send(msg2)
	return err
}
