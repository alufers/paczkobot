package tghelpers

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FakeUpdateHook struct {
	didRun bool
}

func (h *FakeUpdateHook) OnUpdate(ctx context.Context) context.Context {
	h.didRun = true
	return ctx
}

func (h *FakeUpdateHook) OnAfterUpdate(ctx context.Context) context.Context {
	return ctx
}

// a test that chcks if command dispatcher executes update hooks
func TestCommandDispatcherUpdateHooks(t *testing.T) {
	botAPI := &MockBotApi{}
	// create a fake command dispatcher
	dispatcher := NewCommandDispatcher(botAPI, nil)

	fh := &FakeUpdateHook{}
	// register a fake update hook
	dispatcher.RegisterUpdateHooks(fh)

	dispatcher.processIncomingUpdate(context.Background(), tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "test",
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: 456,
			},
		},
	})

	// check if the hook was run
	if !fh.didRun {
		t.Fatal("update hook was not run")
	}
}
