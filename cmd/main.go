package main

import (
	"flag"
	"log"

	"github.com/bells307/poll-telegram-bot/internal/app"
	"github.com/bells307/poll-telegram-bot/internal/app/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	mode      string
	yaml_file string
)

func init() {
	flag.StringVar(&mode, "m", "yaml", "Bot mode")
	flag.StringVar(&yaml_file, "yaml", "opts.yaml", "Yaml file path (for yaml mode)")
	flag.Parse()
}

func main() {
	// Подключаем api
	api, err := tgbotapi.NewBotAPI("2049961870:AAGO5zAwd4aMUdCirLwjn0-05Bjn252bXoU")
	if err != nil {
		log.Panic(err)
	}

	api.Debug = true

	log.Printf("Authorized on account %s", api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Устанавливаем режим работы бота
	var cfg config.Config = nil
	if mode == "yaml" {
		cfg, err = config.NewYamlConfigProvider("opts.yaml")
		if err != nil {
			log.Panic(err)
		}
	}

	// Запуск
	bot := app.NewPollBot(api, cfg)
	bot.Run(updateConfig)
}
