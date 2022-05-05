package paczkobot

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AskService struct {
	BotApp            *BotApp
	AskCallbacks      map[int64]func(string, error)
	AskCallbacksMutex sync.Mutex
}

func NewAskService(botApp *BotApp) *AskService {
	return &AskService{
		BotApp:       botApp,
		AskCallbacks: map[int64]func(string, error){},
	}
}

func (a *AskService) ProcessIncomingMessage(update tgbotapi.Update) bool {
	a.AskCallbacksMutex.Lock()
	defer a.AskCallbacksMutex.Unlock()

	if update.CallbackQuery != nil {
		if update.CallbackQuery.Data == "/cancel" {
			if callback, ok := a.AskCallbacks[update.CallbackQuery.From.ID]; ok {
				a.BotApp.Bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Canceled"))
				callback("", errors.New("canceled"))
				delete(a.AskCallbacks, update.CallbackQuery.From.ID)
			}
			return true
		}
		if update.CallbackQuery.Data == "/yes" {
			if callback, ok := a.AskCallbacks[update.CallbackQuery.From.ID]; ok {
				a.BotApp.Bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Confirmed"))
				callback("", nil)
				delete(a.AskCallbacks, update.CallbackQuery.From.ID)
			}
			return true
		}
		if strings.HasPrefix(update.CallbackQuery.Data, "/sugg ") {
			val := strings.TrimPrefix(update.CallbackQuery.Data, "/sugg ")
			if callback, ok := a.AskCallbacks[update.CallbackQuery.From.ID]; ok {
				a.BotApp.Bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Suggested "+val))
				callback(val, nil)
				delete(a.AskCallbacks, update.CallbackQuery.From.ID)
			}
		}
	}

	if update.Message != nil {
		log.Printf("Processing msg from %v: %v", update.Message.From.ID, update.Message.Text)
		if strings.HasPrefix(update.Message.Text, "/") {
			if callback, ok := a.AskCallbacks[update.Message.From.ID]; ok {
				callback("", errors.New("canceled"))
				delete(a.AskCallbacks, update.Message.From.ID)
				return false
			}
		}

		if callback, ok := a.AskCallbacks[update.Message.From.ID]; ok {
			a.BotApp.Bot.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
			callback(update.Message.Text, nil)
			delete(a.AskCallbacks, update.Message.From.ID)
			return true
		}
	}

	return false
}

func (a *AskService) AskForArgument(chatID int64, question string, suggestionsArr ...map[string]string) (string, error) {

	suggestions := map[string]string{}
	if len(suggestionsArr) != 0 {
		suggestions = suggestionsArr[0]
	}
	extraButtons := [][]tgbotapi.InlineKeyboardButton{}
	extraButtons = append(extraButtons, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "/cancel"),
	})
	if len(suggestions) > 0 {
		for key, value := range suggestions {
			extraButtons = append(extraButtons,
				[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(value, "/sugg "+key)},
			)
		}
	}
	msg := tgbotapi.NewMessage(chatID, question)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: extraButtons,
	}

	msg.ReplyToMessageID = 0
	msg.ParseMode = "HTML"
	sentMsg, err := a.BotApp.Bot.Send(msg)
	if err != nil {
		return "", err
	}

	retChan := make(chan interface{})
	func() {
		a.AskCallbacksMutex.Lock()
		defer a.AskCallbacksMutex.Unlock()
		a.AskCallbacks[chatID] = func(answer string, err error) {
			if err != nil {
				select {
				case retChan <- err:
				default:
				}
				return
			}
			select {
			case retChan <- answer:
			default:

			}
		}
	}()

	timeout := time.After(time.Second * 60 * 10)
	select {
	case answer := <-retChan:
		switch v := answer.(type) {
		case string:
			_, err := a.BotApp.Bot.Send(tgbotapi.NewEditMessageTextAndMarkup(
				chatID,
				sentMsg.MessageID,
				question+" "+v,
				tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
				},
			))
			if err != nil {
				return "", err
			}
			return v, nil
		case error:
			a.BotApp.Bot.Send(tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID))
			return "", v
		default:
			a.BotApp.Bot.Send(tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID))
			return "", errors.New("unknown answer type")
		}
	case <-timeout:
		a.AskCallbacksMutex.Lock()
		defer a.AskCallbacksMutex.Unlock()
		delete(a.AskCallbacks, chatID)
		a.BotApp.Bot.Send(tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID))
		return "", errors.New("timed out while waiting for answer")
	}
}

func (a *AskService) Confirm(chatID int64, question string) error {
	msg := tgbotapi.NewMessage(chatID, question)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("✅ Yes", "/yes"),
				tgbotapi.NewInlineKeyboardButtonData("❌ No", "/no"),
			},
		},
	}

	msg.ReplyToMessageID = 0
	msg.ParseMode = "HTML"
	sentMsg, err := a.BotApp.Bot.Send(msg)
	if err != nil {
		return err
	}
	retChan := make(chan interface{})
	func() {
		a.AskCallbacksMutex.Lock()
		defer a.AskCallbacksMutex.Unlock()
		a.AskCallbacks[chatID] = func(answer string, err error) {

			if err != nil {
				select {
				case retChan <- err:
				default:
				}
				return
			}
			select {
			case retChan <- answer:
			default:

			}
		}
	}()
	defer a.BotApp.Bot.Send(tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID))
	timeout := time.After(time.Second * 60 * 10)
	select {
	case answer := <-retChan:
		switch v := answer.(type) {
		case string:
			return nil
		case error:
			return v
		default:
			return errors.New("unknown answer type")
		}
	case <-timeout:
		a.AskCallbacksMutex.Lock()
		defer a.AskCallbacksMutex.Unlock()
		delete(a.AskCallbacks, chatID)
		return errors.New("timed out while waiting for answer")
	}

}
