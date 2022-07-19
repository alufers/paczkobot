package paczkobot

import (
	"context"
	"fmt"
	"time"
)

type FollowService struct {
	App *BotApp
}

func NewFollowService(app *BotApp) *FollowService {
	return &FollowService{App: app}
}

func (f *FollowService) FollowPackage(ctx context.Context, shipmentNumber string, telegramUserID int64, chatID int64, providersToFollow []*FollowedPackageProvider, followedPackage *FollowedPackage) error {

	followedPackage.TrackingNumber = shipmentNumber
	followedPackage.LastAutomaticCheck = time.Now()
	followedPackage.LastChange = time.Now()
	followedPackage.FollowedPackageProviders = providersToFollow

	if err := f.App.DB.Unscoped().Where("tracking_number = ?", shipmentNumber).FirstOrCreate(followedPackage).Error; err != nil {
		return fmt.Errorf("failed to create FollowedPackage: %v", err)
	}

	if followedPackage.DeletedAt.Valid {
		followedPackage.DeletedAt.Valid = false
		if err := f.App.DB.Save(followedPackage).Error; err != nil {
			return fmt.Errorf("failed to restore FollowedPackage: %v", err)
		}
	}

	followedPackageTelegramUser := &FollowedPackageTelegramUser{
		FollowedPackageID: followedPackage.ID,
		TelegramUserID:    telegramUserID,
		ChatID:            chatID,
	}

	if err := f.App.DB.Where("followed_package_id = ? AND chat_id = ?",
		followedPackage.ID,
		followedPackageTelegramUser.ChatID,
	).FirstOrCreate(followedPackageTelegramUser).Error; err != nil {
		return fmt.Errorf("failed to create FollowedPackageTelegramUser: %v", err)
	}

	for _, p := range providersToFollow {
		p.FollowedPackageID = followedPackage.ID
		if err := f.App.DB.Where("followed_package_id = ? AND provider_name = ?",
			followedPackage.ID,
			p.ProviderName,
		).FirstOrCreate(p).Error; err != nil {
			return fmt.Errorf("failed to create FollowedPackageProvider: %v", err)
		}
	}

	return nil
}
