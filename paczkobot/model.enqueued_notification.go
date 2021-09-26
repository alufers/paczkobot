package paczkobot

type EnqueuedNotification struct {
	Model

	FollowedPackageTelegramUserID string
	FollowedPackageTelegramUser   *FollowedPackageTelegramUser
	TelegramUserID                int // used for querying
	FollowedPackageProviderID     string
	FollowedPackageProvider       *FollowedPackageProvider
}
