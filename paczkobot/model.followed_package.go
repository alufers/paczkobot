package paczkobot

import "time"

type FollowedPackage struct {
	Model
	TrackingNumber               string `gorm:"unique"`
	FollowedPackageProviders     []*FollowedPackageProvider
	FollowedPackageTelegramUsers []*FollowedPackageTelegramUser
	LastAutomaticCheck           time.Time
	LastChange                   time.Time
	Inactive                     bool
}

type FollowedPackageTelegramUser struct {
	Model
	FollowedPackageID string
	FollowedPackage   *FollowedPackage
	TelegramUserID    int
	ChatID            int64
	CustomName        string
}

type FollowedPackageProvider struct {
	Model
	FollowedPackage    *FollowedPackage
	FollowedPackageID  string
	ProviderName       string
	LastStatusDate     time.Time
	LastStatusValue    string
	LastStatusLocation string
}
