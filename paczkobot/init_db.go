package paczkobot

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dbType := viper.GetString("db.type")
	switch dbType {
	case "sqlite":
		return gorm.Open(sqlite.Open(viper.GetString("db.filename")), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	case "postgres":
		return gorm.Open(postgres.Open(viper.GetString("db.dsn")), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unknown database type '%v'", dbType)
	}

}
