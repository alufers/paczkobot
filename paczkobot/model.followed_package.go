package paczkobot

import "time"

type FollowedPackage struct {
	Model
	TrackingNumber               string `unique`
	FollowedPackageProviders     []*FollowedPackageProvider
	FollowedPackageTelegramUsers []*FollowedPackageTelegramUser
	LastCheck                    time.Time
}

type FollowedPackageTelegramUser struct {
	Model

	FollowedPackageID string
	FollowedPackage   *FollowedPackage
	TelegramUserID    int
}

type FollowedPackageProvider struct {
	Model
	FollowedPackage   *FollowedPackage
	FollowedPackageID string
	ProviderName      string
	LastStatusDate    time.Time
	LastStatusValue   string
}
