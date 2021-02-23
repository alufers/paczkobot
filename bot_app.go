package main

import (
	"errors"
	"fmt"
	"github.com/alufers/paczkobot/commondata"
	"html"
	"log"
	"strings"

	"github.com/alufers/paczkobot/commonerrors"
	"github.com/alufers/paczkobot/providers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotApp struct {
	Bot *tgbotapi.BotAPI
}

func NewBotApp(b *tgbotapi.BotAPI) *BotApp {
	return &BotApp{
		Bot: b,
	}
}

func (a *BotApp) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := a.Bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatalf("telegram updates error: %v", err)
	}
	log.Printf("Telegram bot is starting...")

	for u := range updates {
		go func(update tgbotapi.Update) {
			var err error
			log.Printf("msg: %v", update.Message.Text)
			if strings.HasPrefix(update.Message.Text, "/track") {
				err = a.handleTrackCommand(update)
			}
			log.Print(err)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Error: <b>"+html.EscapeString(err.Error())+"</b>")
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID

				a.Bot.Send(msg)
			}
		}(u)
	}

}

type providerReply struct {
	provider providers.Provider
	data     *commondata.TrackingData
	err      error
}

func (a *BotApp) handleTrackCommand(update tgbotapi.Update) error {
	var segments = strings.Split(update.Message.Text, " ")
	log.Printf("segments = %#v", segments)
	if len(segments) < 2 {
		return fmt.Errorf("usage: /track &lt;shipmentNumber&gt;")
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

	statuses := map[string]string{}
	replyChan := make(chan *providerReply, len(providersToCheck))
	for _, p := range providersToCheck {
		statuses[p.GetName()] = "⌛ checking..."
		go func(p providers.Provider) {
			d, err := p.Track(shipmentNumber)
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
		for n, v := range statuses {
			msgText += fmt.Sprintf("%v: <b>%v</b>\n", n, html.EscapeString(v))
		}
		if msgIdToEdit != 0 {
			msg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, msgIdToEdit, msgText)
			msg.ParseMode = "HTML"
			_, err := a.Bot.Send(msg)
			if err != nil {
				log.Printf("failed to edit status msg: %v", err)
				return
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ParseMode = "HTML"
			// msg.ReplyToMessageID = update.Message.MessageID

			res, err := a.Bot.Send(msg)
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

			var longTracking = fmt.Sprintf("Detailed tracking for package <i>%v</i> provided by <i>%v</i>:\n", rep.data.ShipmentNumber, rep.data.ProviderName)

			for i, ts := range rep.data.TrackingSteps {
				shouldBold := i == len(rep.data.TrackingSteps)-1
				if shouldBold {
					longTracking += "<b>"
				}
				longTracking += ts.Datetime.Format("2006-01-02 15:04") + " " + ts.Message
				if ts.Location != "" {
					longTracking += " 📌 " + ts.Location
				}
				longTracking += "\n"
				if shouldBold {
					longTracking += "</b>"
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, longTracking)
			msg.ParseMode = "HTML"
			msg.ReplyToMessageID = update.Message.MessageID

			a.Bot.Send(msg)
		}
	}

	return nil
}
