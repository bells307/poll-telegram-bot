package main

import (
	"flag"
	"log"
	"os"

	"github.com/bells307/poll-telegram-bot/internal/app"
	"github.com/bells307/poll-telegram-bot/internal/app/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	token     string
	mode      string
	yaml_file string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot token")
	flag.StringVar(&mode, "m", "yaml", "Bot mode")
	flag.StringVar(&yaml_file, "f", "config.yaml", "Yaml file path (for yaml mode)")
	flag.Parse()
}

func main() {
	if len(token) == 0 {
		log.Panic("token is empty")
	}

	// Подключаем api
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	if os.Getenv("BOT_DEBUG") == "true" {
		api.Debug = true
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Устанавливаем режим работы бота
	var cfg config.Config = nil
	if mode == "yaml" {
		cfg, err = config.NewYamlConfigProvider(yaml_file)
		if err != nil {
			log.Panic(err)
		}
	}

	// Запуск
	bot := app.NewPollBot(api, cfg)
	bot.Run(updateConfig)
}
