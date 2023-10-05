package paczkobot

import (
	"context"
	"errors"
	"fmt"
	"html"
	"log"
	"sort"

	"github.com/alufers/paczkobot/commondata"
	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/providers"
	"github.com/alufers/paczkobot/tghelpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TrackCommand struct {
	App *BotApp
}

func (s *TrackCommand) Aliases() []string {
	return []string{"/track"}
}

func (s *TrackCommand) Arguments() []*tghelpers.CommandDefArgument {
	return []*tghelpers.CommandDefArgument{
		{
			Name:        "shipmentNumber",
			Description: "shipment number of the package",
			Question:    "Please enter the shipment number to track:",
		},
	}
}

func (s TrackCommand) Category() string {
	return "Tracking"
}

func (t *TrackCommand) Help() string {
	return "shows up-to-date tracking information about a package with the given number"
}

type providerReply struct {
	provider providers.Provider
	data     *commondata.TrackingData
	err      error
}

func (t *TrackCommand) Execute(ctx context.Context) error {
	args := tghelpers.ArgsFromCtx(ctx)
	shipmentNumber, err := args.GetOrAskForArgument("shipmentNumber")
	if err != nil {
		return err
	}
	providersToCheck := []providers.Provider{}
	for _, provider := range providers.AllProviders {
		if provider.MatchesNumber(shipmentNumber) {
			providersToCheck = append(providersToCheck, provider)
		}
	}

	if len(providersToCheck) == 0 {
		return fmt.Errorf("no tracking providers support this tracking number")
	}

	statuses := map[string]string{}
	replyChan := make(chan *providerReply, len(providersToCheck))
	for _, p := range providersToCheck {
		statuses[p.GetName()] = "⌛ checking..."
		go func(p providers.Provider) {
			d, err := t.App.TrackingService.InvokeProviderAndNotifyFollowers(ctx, p, shipmentNumber)
			if err != nil {
				replyChan <- &providerReply{
					provider: p,
					err:      err,
				}
			} else {
				replyChan <- &providerReply{
					provider: p,
					data:     d,
				}
			}
		}(p)
	}
	var msgIdToEdit int
	sendStatuses := func() {
		var msgText string
		statusesKeys := []string{}
		for k := range statuses {
			statusesKeys = append(statusesKeys, k)
		}
		sort.Strings(statusesKeys)
		for _, n := range statusesKeys {
			v := statuses[n]
			msgText += fmt.Sprintf("%v: <b>%v</b>\n", n, html.EscapeString(v))
		}
		if msgIdToEdit != 0 {
			msg := tgbotapi.NewEditMessageText(args.ChatID, msgIdToEdit, msgText)
			msg.ParseMode = "HTML"
			_, err := t.App.Bot.Send(msg)
			if err != nil {
				log.Printf("failed to edit status msg: %v", err)
				return
			}
		} else {
			msg := tgbotapi.NewMessage(args.ChatID, msgText)
			msg.ParseMode = "HTML"
			// msg.ReplyToMessageID = update.Message.MessageID

			res, err := t.App.Bot.Send(msg)
			if err != nil {
				log.Printf("failed to send status msg: %v", err)
				return
			}
			msgIdToEdit = res.MessageID
		}
	}
	sendStatuses()
	for rep := range replyChan {
		if rep.err != nil {
			if errors.Is(rep.err, commonerrors.NotFoundError) {
				statuses[rep.provider.GetName()] = "🔳 Not found"
			} else {
				statuses[rep.provider.GetName()] = "⚠️ Error: " + rep.err.Error()
			}

			sendStatuses()
		} else {
			status := ""
			if len(rep.data.TrackingSteps) > 0 {
				status = rep.data.TrackingSteps[len(rep.data.TrackingSteps)-1].Message
			}
			statuses[rep.provider.GetName()] = "🔎 " + status
			sendStatuses()

			longTracking := fmt.Sprintf("Detailed tracking for package <i>%v</i> provided by <b>%v</b>:\n", rep.data.ShipmentNumber, rep.data.ProviderName)

			for i, ts := range rep.data.TrackingSteps {
				shouldBold := i == len(rep.data.TrackingSteps)-1
				if shouldBold {
					longTracking += "<b>"
				}
				emoji := commondata.CommonTrackingStepTypeEmoji[ts.CommonType]
				if emoji != "" {
					emoji += " "
				}
				longTracking += ts.Datetime.Format("2006-01-02 15:04") + " " + emoji + ts.Message
				if ts.Location != "" {
					longTracking += " 📌 " + ts.Location
				}
				longTracking += "\n"
				if shouldBold {
					longTracking += "</b>"
				}
			}

			detailsString := ""
			if rep.data.Destination != "" && rep.data.SentFrom != "" {
				detailsString += "The package is headed from <u>" + rep.data.SentFrom + "</u> to <u>" + rep.data.Destination + "</u>."
			} else if rep.data.Destination != "" {
				detailsString += "The package is headed to <u>" + rep.data.Destination + "</u>."
			}

			if rep.data.Weight != 0.0 {
				if detailsString != "" {
					detailsString += " "
				}
				detailsString += fmt.Sprintf("The package weighs %.2f kg.", rep.data.Weight)
			}

			if detailsString != "" {
				longTracking += "\n" + detailsString + "\n"
			}

			msg := tgbotapi.NewMessage(args.ChatID, longTracking)
			msg.ParseMode = "HTML"
			if args.Update.Message != nil {
				msg.ReplyToMessageID = args.Update.Message.MessageID
			}
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🚶 Follow this package", fmt.Sprintf("/follow %v", rep.data.ShipmentNumber)),
				),
			)
			_, err := t.App.Bot.Send(msg)
			if err != nil {
				return err
			}

		}
	}

	return nil
}
