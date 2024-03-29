package paczkobot

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alufers/paczkobot/providers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type NotificationsService struct {
	app   *BotApp
	Hooks []PackageNotificationHook
}

func NewNotificationsService(app *BotApp, hooks []PackageNotificationHook) *NotificationsService {
	return &NotificationsService{
		app:   app,
		Hooks: hooks,
	}
}

// NotifyProviderStatusChanged is executed from the tracking service whenever a status reported by a provider is changed.
func (s *NotificationsService) NotifyProviderStatusChanged(provider providers.Provider, followedPackage *FollowedPackage) error {
	var prov *FollowedPackageProvider
	for _, p := range followedPackage.FollowedPackageProviders {
		if p.ProviderName == provider.GetName() {
			prov = p
			break
		}
	}

	if prov == nil {
		return fmt.Errorf("provider %v not found in package %v while enqueuing notification", provider.GetName(), followedPackage.TrackingNumber)
	}

	for _, u := range followedPackage.FollowedPackageTelegramUsers {
		notif := &EnqueuedNotification{
			FollowedPackageTelegramUser: u,
			TelegramUserID:              u.TelegramUserID,
			ChatID:                      u.ChatID,
			FollowedPackageProvider:     prov,
		}
		if err := s.app.DB.Save(notif).Error; err != nil {
			return fmt.Errorf("failed to enqueue notification for TG user %v, chat ID %v, provider %v: %w", u.TelegramUserID, u.ChatID, provider.GetName(), err)
		}
		go func(u *FollowedPackageTelegramUser) {
			time.Sleep(time.Second * 30)
			err := s.sendNotificationsForUser(u)
			if err != nil {
				log.Printf("failed to send notifications for user %v: %v", u.TelegramUserID, err)
			}
		}(u)
	}

	return nil
}

func (s *NotificationsService) sendNotificationsForUser(tgUser *FollowedPackageTelegramUser) error {
	var notifications []*EnqueuedNotification
	if err := s.app.DB.Where("chat_id = ?", tgUser.ChatID).
		Preload("FollowedPackageProvider").
		Preload("FollowedPackageTelegramUser").
		Preload("FollowedPackageProvider.FollowedPackage").
		Find(&notifications).Error; err != nil {
		return fmt.Errorf("failed to fetch notifications for user %v: %w", tgUser.TelegramUserID, err)
	}
	if len(notifications) == 0 {
		return nil
	}
	msgContents := ""

	for _, n := range notifications {
		loc := ""
		if strings.TrimSpace(n.FollowedPackageProvider.LastStatusLocation) != "" {
			loc = " 📌" + n.FollowedPackageProvider.LastStatusLocation
		}
		customName := ""
		if n.FollowedPackageTelegramUser.CustomName != "" {
			customName = n.FollowedPackageTelegramUser.CustomName + " "
		}
		msgContents += fmt.Sprintf("%v(%v) <b>%v</b>: %v %v%v\n",
			customName,
			n.FollowedPackageProvider.ProviderName,
			n.FollowedPackageProvider.FollowedPackage.TrackingNumber,
			n.FollowedPackageProvider.LastStatusValue,
			n.FollowedPackageProvider.LastStatusDate.Format("2006-01-02 15:04"),
			loc,
		)
		if err := s.app.DB.Delete(n).Error; err != nil {
			return fmt.Errorf("failed to delete notification %v: %w", n.ID, err)
		}

	}

	msg := tgbotapi.NewMessage(notifications[0].FollowedPackageTelegramUser.ChatID, msgContents)
	msg.ParseMode = "HTML"
	keyboard := [][]tgbotapi.InlineKeyboardButton{}
	for _, hook := range s.Hooks {
		for _, n := range notifications {
			hookResult, err := hook.HookNotificationKeyboard(n)
			if err != nil {
				log.Printf("failed to get keyboard from hook %T: %v", hook, err)
				continue
			}
			if hookResult != nil {
				keyboard = append(keyboard, hookResult...)
			}
		}
	}
	if len(keyboard) > 0 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	}

	if _, err := s.app.Bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send notifications to chat %v: %w", tgUser.ChatID, err)
	}

	return nil
}

func (s *NotificationsService) FlushEnqueuedNotifications() error {
	limit := 0
	for {
		limit++
		if limit > 59 {
			log.Printf("Bailing after sending 59 enqueued notifications")
			return nil
		}
		var nextNotification *EnqueuedNotification
		result := s.app.DB.Preload("FollowedPackageTelegramUser").Limit(1).Find(&nextNotification)
		if err := result.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed to fetch notifications to flush: %w", err)
		}

		if nextNotification == nil || result.RowsAffected < 1 {
			return nil
		}

		if err := s.sendNotificationsForUser(nextNotification.FollowedPackageTelegramUser); err != nil {
			return err
		}

	}
}
