package paczkobot

import (
	"log"
	"time"

	"github.com/alufers/paczkobot/providers"
	"github.com/alufers/paczkobot/providers/mock"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

func Run() {
	viper.SetConfigName("paczkobot")                       // name of config file (without extension)
	viper.SetConfigType("yaml")                            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                               // optionally look for config in the working directory
	viper.AddConfigPath("/etc/paczkobot/")                 // path to look for the config file in
	viper.AddConfigPath("$HOME/.config/alufers/paczkobot") // call multiple times to add many search paths
	viper.SetDefault("telegram.token", "")
	viper.SetDefault("telegram.username", "paczko_bot")
	viper.SetDefault("telegram.debug", false)

	viper.SetDefault("tracking.max_time_without_change", time.Hour*24*14)
	viper.SetDefault("tracking.automatic_tracking_check_interval", time.Minute*20)
	viper.SetDefault("tracking.automatic_tracking_check_jitter", time.Minute*7)
	viper.SetDefault("tracking.max_packages_per_automatic_tracking_check", 15)
	viper.SetDefault("tracking.delay_between_packages_in_automatic_tracking", time.Minute)
	if viper.GetBool("tracking.mock_provider") {
		providers.AllProviders = append(providers.AllProviders, &mock.MockProvider{})
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		viper.SafeWriteConfig()
		log.Fatalf("config file: %v", err)
	}

	db, err := InitDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", db)
	}

	if err := db.AutoMigrate(&FollowedPackage{}, &FollowedPackageProvider{}, &FollowedPackageTelegramUser{}, &EnqueuedNotification{}); err != nil {
		log.Fatalf("failed to AutoMigrate: %v", err)
	}

	token := viper.GetString("telegram.token")
	if token == "" {
		log.Fatal("no token provided")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = viper.GetBool("telegram.debug")
	app := NewBotApp(bot, db)
	app.Run()
}
