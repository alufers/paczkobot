package tghelpers

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AskService allows commands to interactively ask questions to the user.
// All question functions block the calling goroutine until an answer is given or a timeout happens
type AskService struct {
	AskCallbacks      map[int64]func(string, error)
	AskCallbacksMutex sync.Mutex
	bot               BotAPI
}

func NewAskService(bot BotAPI) *AskService {
	return &AskService{
		bot:          bot,
		AskCallbacks: map[int64]func(string, error){},
	}
}

// Implements UpdateHook
func (a *AskService) OnUpdate(ctx context.Context, update tgbotapi.Update) bool {
	a.AskCallbacksMutex.Lock()
	defer a.AskCallbacksMutex.Unlock()
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Message == nil || update.CallbackQuery.Message.Chat == nil {
			return false
		}
		chatID := update.CallbackQuery.Message.Chat.ID
		if update.CallbackQuery.Data == "/cancel" {
			if callback, ok := a.AskCallbacks[chatID]; ok {
				_, err := a.bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Canceled"))
				if err != nil {
					log.Printf("Error sending callback: %v", err)
				}
				callback("", errors.New("canceled"))
				delete(a.AskCallbacks, chatID)
			}
			return true
		}
		if update.CallbackQuery.Data == "/yes" {
			if callback, ok := a.AskCallbacks[chatID]; ok {
				_, err := a.bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Confirmed"))
				if err != nil {
					log.Printf("Error sending callback: %v", err)
				}
				callback("", nil)
				delete(a.AskCallbacks, chatID)
			}
			return true
		}
		if strings.HasPrefix(update.CallbackQuery.Data, "/sugg ") {
			val := strings.TrimPrefix(update.CallbackQuery.Data, "/sugg ")
			if callback, ok := a.AskCallbacks[chatID]; ok {
				_, err := a.bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Suggested "+val))
				if err != nil {
					log.Printf("Error sending callback: %v", err)
				}
				callback(val, nil)
				delete(a.AskCallbacks, chatID)
			}
		}
	}

	if update.Message != nil {
		log.Printf("Processing msg from chat ID %v: %v", update.Message.Chat.ID, update.Message.Text)
		if strings.HasPrefix(update.Message.Text, "/") {
			if callback, ok := a.AskCallbacks[update.Message.Chat.ID]; ok {
				callback("", errors.New("canceled"))
				delete(a.AskCallbacks, update.Message.Chat.ID)
				return false
			}
		}

		if callback, ok := a.AskCallbacks[update.Message.Chat.ID]; ok {
			_, err := a.bot.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
			if err != nil {
				log.Printf("Error deleting message: %v", err)
			}
			callback(update.Message.Text, nil)
			delete(a.AskCallbacks, update.Message.Chat.ID)
			return true
		}
	}

	return false
}

// AskForArgument asks the user at the specified chatID for a text value.
// suggestionsArr should contain the map of suggestions where the key is the value that will be returned and the value is the text that will be displayed to the user.
func (a *AskService) AskForArgument(chatID int64, question string, suggestionsArr ...map[string]string) (string, error) {
	suggestions := map[string]string{}
	if len(suggestionsArr) != 0 {
		suggestions = suggestionsArr[0]
	}
	extraButtons := [][]tgbotapi.InlineKeyboardButton{}
	extraButtons = append(extraButtons, []tgbotapi.InlineKeyboardButton{
		// tgbotapi.NewInlineKeyboardButtonData
		tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "/cancel"),
	})
	if len(suggestions) > 0 {
		for key, value := range suggestions {
			extraButtons = append(extraButtons,
				[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(value, "/sugg "+key)},
			)
		}
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

	msg := tgbotapi.NewMessage(chatID, question)
	if len(suggestions) > 0 {
		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: extraButtons,
		}
	} else {
		msg.ReplyMarkup = &tgbotapi.ForceReply{
			ForceReply:            true,
			InputFieldPlaceholder: question,
		}
	}

	msg.ReplyToMessageID = 0
	msg.ParseMode = "HTML"
	_, err := a.bot.Send(msg) // sendMsg
	if err != nil {
		return "", err
	}

	timeout := time.After(time.Second * 60 * 10)
	select {
	case answer := <-retChan:
		switch v := answer.(type) {
		case string:
			// _, err := a.bot.Send(tgbotapi.NewEditMessageText(
			// 	chatID,
			// 	sentMsg.MessageID,
			// 	question+" "+v,
			// 	// tgbotapi.InlineKeyboardMarkup{
			// 	// 	InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
			// 	// },
			// ))
			// if err != nil {
			// 	return "", fmt.Errorf("failed to edit question message: %w", err)
			// }
			return v, nil
		case error:
			// a.bot.Send(tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID))
			return "", v
		default:
			// a.bot.Send(tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID))
			return "", errors.New("unknown answer type")
		}
	case <-timeout:
		a.AskCallbacksMutex.Lock()
		defer a.AskCallbacksMutex.Unlock()
		delete(a.AskCallbacks, chatID)
		// a.bot.Send(tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID))
		return "", errors.New("timed out while waiting for answer")
	}
}

func (a *AskService) Confirm(chatID int64, question string) error {
	msg := tgbotapi.NewMessage(chatID, question)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("✅ Yes", "/yes"),
				tgbotapi.NewInlineKeyboardButtonData("❌ No", "/cancel"),
			},
		},
	}

	msg.ReplyToMessageID = 0
	msg.ParseMode = "HTML"
	_, err := a.bot.Send(msg) // sentMsg
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
	// defer a.bot.Send(tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID))
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
