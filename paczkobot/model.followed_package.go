package paczkobot

import (
	"time"

	"github.com/alufers/paczkobot/dbutil"
	"github.com/alufers/paczkobot/inpostextra"
)

type FollowedPackage struct {
	dbutil.Model
	TrackingNumber               string `gorm:"unique"`
	FollowedPackageProviders     []*FollowedPackageProvider
	FollowedPackageTelegramUsers []*FollowedPackageTelegramUser
	LastAutomaticCheck           time.Time
	LastChange                   time.Time
	Inactive                     bool
	FromName                     string
	InpostCredentialsID          string
	InpostCredentials            *inpostextra.InpostCredentials
}

type FollowedPackageTelegramUser struct {
	dbutil.Model
	FollowedPackageID string
	FollowedPackage   *FollowedPackage
	TelegramUserID    int64
	ChatID            int64
	CustomName        string
}

type FollowedPackageProvider struct {
	dbutil.Model
	FollowedPackage    *FollowedPackage
	FollowedPackageID  string
	ProviderName       string
	LastStatusDate     time.Time
	LastStatusValue    string
	LastStatusLocation string
}
