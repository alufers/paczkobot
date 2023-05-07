package paczkobot

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/providers"
	"github.com/alufers/paczkobot/tghelpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FollowCommand struct {
	App *BotApp
}

func (s *FollowCommand) Aliases() []string {
	return []string{"/follow"}
}

func (s *FollowCommand) Arguments() []*tghelpers.CommandDefArgument {
	return []*tghelpers.CommandDefArgument{
		{
			Name:        "shipmentNumber",
			Description: "shipment number of the package",
			Question:    "Please enter the shipment number:",
		},
	}
}

func (s *FollowCommand) Category() string {
	return "Following packages"
}

func (f *FollowCommand) Help() string {
	return "follows a package and sends you an update every time its status changes"
}

func (f *FollowCommand) Execute(ctx context.Context, args *tghelpers.CommandArguments) error {
	msg := tgbotapi.NewMessage(args.ChatID, "⌛ loading...")
	msg.ParseMode = "HTML"
	loadingRes, err := f.App.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send loading message: %w", err)
	}
	defer func() {
		_, err := f.App.Bot.Send(tgbotapi.NewDeleteMessage(args.ChatID, loadingRes.MessageID))
		if err != nil {
			log.Printf("failed to delete loading message: %v", err)
		}
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

	err = f.App.FollowService.FollowPackage(
		ctx,
		shipmentNumber,
		args.FromUserID,
		args.ChatID,
		providersToFollow,
		&FollowedPackage{},
	)
	if err != nil {
		return fmt.Errorf("failed to follow package: %w", err)
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
