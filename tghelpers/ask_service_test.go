package tghelpers_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/alufers/paczkobot/tghelpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestAskServiceReturnsFalseForUnrelatedUpdates(t *testing.T) {
	botApi := &tghelpers.MockBotApi{}
	askService := tghelpers.NewAskService(botApi)
	res := askService.OnUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "foo",
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
		},
	})
	assert.False(t, res)
}

func TestAskServiceConfirmWorks(t *testing.T) {
	stage := 0
	failChan := make(chan error, 1)
	var askService *tghelpers.AskService
	botApi := &tghelpers.MockBotApi{
		SendCallback: func(c tgbotapi.Chattable) (tgbotapi.Message, error) {
			switch stage {
			case 0:
				assert.Equal(t, "Test question?", c.(tgbotapi.MessageConfig).Text)

				stage++
				assert.IsType(t, tgbotapi.InlineKeyboardMarkup{}, c.(tgbotapi.MessageConfig).ReplyMarkup)
				assert.Equal(t, 2, len(c.(tgbotapi.MessageConfig).ReplyMarkup.(tgbotapi.InlineKeyboardMarkup).InlineKeyboard[0]))
				assert.Contains(t, c.(tgbotapi.MessageConfig).ReplyMarkup.(tgbotapi.InlineKeyboardMarkup).InlineKeyboard[0][0].Text, "Yes")
				msg := tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: 1234,
					},
				}
				go func() {
					time.Sleep(1 * time.Millisecond)
					res := askService.OnUpdate(tgbotapi.Update{
						CallbackQuery: &tgbotapi.CallbackQuery{
							ID:      "123",
							Message: &msg,
							Data:    *c.(tgbotapi.MessageConfig).ReplyMarkup.(tgbotapi.InlineKeyboardMarkup).InlineKeyboard[0][0].CallbackData,
						},
					})
					assert.True(t, res) // should return true because it's a related update
				}()

				return msg, nil
			case 1:
				assert.IsType(t, tgbotapi.CallbackConfig{}, c)
				stage++
				return tgbotapi.Message{}, nil
			default:
				failChan <- fmt.Errorf("Unexpected Send call")
				t.Fatal("Unexpected Send call")
				return tgbotapi.Message{}, nil
			}
		},
	}
	askService = tghelpers.NewAskService(botApi)
	chatID := int64(1234)
	resultChan := make(chan error)
	go func() {
		resultErr := askService.Confirm(chatID, "Test question?")
		resultChan <- resultErr
	}()

	select {
	case result := <-resultChan:
		assert.NoError(t, result)
		assert.Equal(t, 2, stage)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout")
	case err := <-failChan:
		t.Fatal(err)
	}
}

func TestAskServiceAskForArgumentWorks(t *testing.T) {
	stage := 0
	failChan := make(chan error, 1)
	var askService *tghelpers.AskService
	botApi := &tghelpers.MockBotApi{
		SendCallback: func(c tgbotapi.Chattable) (tgbotapi.Message, error) {
			switch stage {
			case 0:
				assert.Equal(t, "Test question?", c.(tgbotapi.MessageConfig).Text)

				stage++
				assert.IsType(t, tgbotapi.InlineKeyboardMarkup{}, c.(tgbotapi.MessageConfig).ReplyMarkup)
				assert.Equal(t, 2, len(c.(tgbotapi.MessageConfig).ReplyMarkup.(tgbotapi.InlineKeyboardMarkup).InlineKeyboard))
				assert.Contains(t, c.(tgbotapi.MessageConfig).ReplyMarkup.(tgbotapi.InlineKeyboardMarkup).InlineKeyboard[0][0].Text, "Cancel")
				assert.Contains(t, c.(tgbotapi.MessageConfig).ReplyMarkup.(tgbotapi.InlineKeyboardMarkup).InlineKeyboard[1][0].Text, "bar")
				msg := tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: 1234,
					},
				}
				go func() {
					time.Sleep(1 * time.Millisecond)
					askService.OnUpdate(tgbotapi.Update{
						CallbackQuery: &tgbotapi.CallbackQuery{
							ID:      "123",
							Message: &msg,
							Data:    *c.(tgbotapi.MessageConfig).ReplyMarkup.(tgbotapi.InlineKeyboardMarkup).InlineKeyboard[0][0].CallbackData,
						},
					})
				}()

				return msg, nil
			case 1:
				assert.IsType(t, tgbotapi.CallbackConfig{}, c)
				stage++
				return tgbotapi.Message{}, nil
			default:
				failChan <- fmt.Errorf("Unexpected Send call")
				t.Fatal("Unexpected Send call")
				return tgbotapi.Message{}, nil
			}
		},
	}
	askService = tghelpers.NewAskService(botApi)
	chatID := int64(1234)
	resultChan := make(chan error)
	go func() {
		_, resultErr := askService.AskForArgument(chatID, "Test question?", map[string]string{
			"foo": "bar",
		})

		resultChan <- resultErr
	}()

	select {
	case result := <-resultChan:
		assert.Error(t, result)
		assert.Equal(t, 2, stage)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout")
	case err := <-failChan:
		t.Fatal(err)
	}
}
