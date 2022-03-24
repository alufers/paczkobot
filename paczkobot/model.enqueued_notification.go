package paczkobot

import "github.com/alufers/paczkobot/dbutil"

type EnqueuedNotification struct {
	dbutil.Model

	FollowedPackageTelegramUserID string
	FollowedPackageTelegramUser   *FollowedPackageTelegramUser
	TelegramUserID                int64 // used for querying
	FollowedPackageProviderID     string
	FollowedPackageProvider       *FollowedPackageProvider
}
