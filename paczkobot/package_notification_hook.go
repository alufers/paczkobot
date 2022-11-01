package paczkobot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type PackageNotificationHook interface {
	HookNotificationKeyboard(notification *EnqueuedNotification) ([][]tgbotapi.InlineKeyboardButton, error)
}
