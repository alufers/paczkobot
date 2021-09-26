package paczkobot

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/providers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type FollowCommand struct {
	App *BotApp
}

func (f *FollowCommand) Usage() string {
	return "/follow <shipmentNumber>"
}

func (f *FollowCommand) Help() string {
	return "follows a package and sends you an update every time its status changes"
}

func (f *FollowCommand) Execute(ctx context.Context, args *CommandArguments) error {
	msg := tgbotapi.NewMessage(args.update.Message.Chat.ID, "⌛ loading...")
	msg.ParseMode = "HTML"
	loadingRes, err := f.App.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send loading message: %w", err)
	}
	defer func() {
		f.App.Bot.DeleteMessage(tgbotapi.NewDeleteMessage(args.update.Message.Chat.ID, loadingRes.MessageID))
	}()

	var segments = strings.Split(args.update.Message.Text, " ")

	if len(segments) < 2 {
		return fmt.Errorf("usage: /follow &lt;shipmentNumber&gt;")
	}
	shipmentNumber := segments[1]
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
			ProviderName:    p.GetName(),
			LastStatusValue: lastStep.Message,
			LastStatusDate:  lastStep.Datetime,
		})
	}

	if len(providersToFollow) <= 0 {
		return fmt.Errorf("the package was not found in any of the providers, please double check your package number")
	}

	followedPackage := &FollowedPackage{
		TrackingNumber:           shipmentNumber,
		LastCheck:                time.Now(),
		FollowedPackageProviders: providersToFollow,
	}

	if err := f.App.DB.Where("tracking_number = ?", shipmentNumber).FirstOrCreate(followedPackage).Error; err != nil {
		return fmt.Errorf("failed to create FollowedPackage: %v", err)
	}

	followedPackageTelegramUser := &FollowedPackageTelegramUser{
		FollowedPackageID: followedPackage.ID,
		TelegramUserID:    args.update.Message.From.ID,
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
	msg2 := tgbotapi.NewMessage(args.update.Message.Chat.ID, fmt.Sprintf(`Package %v has been added to your followed packages!`))
	msg2.ParseMode = "HTML"
	_, err = f.App.Bot.Send(msg2)
	return err
}
