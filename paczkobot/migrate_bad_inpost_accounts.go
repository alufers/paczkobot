package paczkobot

import (
	"log"

	"gorm.io/gorm"
)

func MigrateBadInpostAccounts(db *gorm.DB) error {
	// Begin a new transaction
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// copy telegram_user_id to telegram_chat_id if telegram_chat_id is 0 in inpost_credentials
	result := tx.Exec(`
		UPDATE inpost_credentials
		SET telegram_chat_id = telegram_user_id
		WHERE telegram_chat_id = 0
	`)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	log.Printf("Updated rows in inpost_credentials: %d", result.RowsAffected)

	// copy telegram_user_id to chat_id if chat_id is 0 in followed_package_telegram_users
	result = tx.Exec(`
		UPDATE followed_package_telegram_users
		SET chat_id = telegram_user_id
		WHERE chat_id = 0
	`)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	log.Printf("Updated rows in followed_package_telegram_users: %d", result.RowsAffected)

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
