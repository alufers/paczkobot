package paczkobot

import (
	"github.com/alufers/paczkobot/commondata"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InpostPackageNotificationHook struct{}

func (h *InpostPackageNotificationHook) HookNotificationKeyboard(notification *EnqueuedNotification) ([][]tgbotapi.InlineKeyboardButton, error) {
	if notification.FollowedPackageProvider.ProviderName != "inpost" {
		return nil, nil
	}
	if notification.FollowedPackageProvider.LastStatusCommonType != commondata.CommonTrackingStepType_READY_FOR_PICKUP {
		return nil, nil
	}
	qrCmd := "/inpostqr " + notification.FollowedPackageProvider.FollowedPackage.TrackingNumber
	openCmd := "/inpostopen " + notification.FollowedPackageProvider.FollowedPackage.TrackingNumber
	return [][]tgbotapi.InlineKeyboardButton{
		{
			{
				Text:         "🔳 QR code",
				CallbackData: &qrCmd,
			},
			{
				Text:         "📦 Open Locker",
				CallbackData: &openCmd,
			},
		},
	}, nil
}
