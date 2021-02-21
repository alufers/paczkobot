package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("paczkobot")                       // name of config file (without extension)
	viper.SetConfigType("yaml")                            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                               // optionally look for config in the working directory
	viper.AddConfigPath("/etc/paczkobot/")                 // path to look for the config file in
	viper.AddConfigPath("$HOME/.config/alufers/paczkobot") // call multiple times to add many search paths
	viper.SetDefault("telegram.token", "")
	viper.SetDefault("telegram.debug", false)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		viper.SafeWriteConfig()
		log.Fatalf("config file: %v", err)
	}
	token := viper.GetString("telegram.token")
	if token == "" {
		log.Fatal("no token provided")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = viper.GetBool("telegram.debug")
	app := NewBotApp(bot)
	app.Run()
}
