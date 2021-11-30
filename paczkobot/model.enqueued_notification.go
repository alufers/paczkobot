package paczkobot

type EnqueuedNotification struct {
	Model

	FollowedPackageTelegramUserID string
	FollowedPackageTelegramUser   *FollowedPackageTelegramUser
	TelegramUserID                int64 // used for querying
	FollowedPackageProviderID     string
	FollowedPackageProvider       *FollowedPackageProvider
}
