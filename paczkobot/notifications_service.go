package paczkobot

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alufers/paczkobot/providers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

type NotificationsService struct {
	app *BotApp
}

func NewNotificationsService(app *BotApp) *NotificationsService {
	return &NotificationsService{
		app: app,
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
			FollowedPackageProvider:     prov,
		}
		if err := s.app.DB.Save(notif).Error; err != nil {
			return fmt.Errorf("failed to enqueue notification for TG user %v, provider %v: %w", u.TelegramUserID, provider.GetName(), err)
		}
		go func(u *FollowedPackageTelegramUser) {
			time.Sleep(time.Second * 30)
			s.sendNotificationsForUser(u)
		}(u)
	}

	return nil
}

func (s *NotificationsService) sendNotificationsForUser(tgUser *FollowedPackageTelegramUser) error {
	var notifications []*EnqueuedNotification
	if err := s.app.DB.Where("telegram_user_id = ?", tgUser.TelegramUserID).
		Preload("FollowedPackageProvider").
		Preload("FollowedPackageTelegramUser").
		Preload("FollowedPackageProvider.FollowedPackage").
		Find(&notifications).Error; err != nil {

		return fmt.Errorf("failed to fetch notifications for user %v: %w", tgUser.TelegramUserID, err)
	}
	if len(notifications) == 0 {
		return nil
	}
	msgContents := "The following packages have new updates: \n"

	for _, n := range notifications {
		loc := ""
		if strings.TrimSpace(n.FollowedPackageProvider.LastStatusLocation) != "" {
			loc = " 📌" + n.FollowedPackageProvider.LastStatusLocation
		}
		msgContents += fmt.Sprintf("<b>%v</b> (%v): %v %v%v\n",
			n.FollowedPackageProvider.FollowedPackage.TrackingNumber,
			n.FollowedPackageProvider.ProviderName,
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
	if _, err := s.app.Bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send notifications to user %v: %w", tgUser.TelegramUserID, err)
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
		if err := s.app.DB.Preload("FollowedPackageTelegramUser").First(&nextNotification).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed to fetch notifications to flush: %w", err)
		}

		if err := s.sendNotificationsForUser(nextNotification.FollowedPackageTelegramUser); err != nil {
			return err
		}

	}
}
