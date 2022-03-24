package inpostextra

import (
	"database/sql"

	"github.com/alufers/paczkobot/dbutil"
)

type InpostCredentials struct {
	dbutil.Model
	TelegramUserID int64
	TelegramChatID int64
	PhoneNumber    string
	AuthToken      string
	RefreshToken   string
	LastScan       sql.NullTime
}
