package main

import (
	"log"

	"github.com/bells307/poll-telegram-bot/internal/app"
	"github.com/bells307/poll-telegram-bot/internal/app/poll_options"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	api, err := tgbotapi.NewBotAPI("2049961870:AAGO5zAwd4aMUdCirLwjn0-05Bjn252bXoU")
	if err != nil {
		log.Panic(err)
	}

	api.Debug = true

	log.Printf("Authorized on account %s", api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	pollOptsProv, err := poll_options.NewFilePollOptions("opts.txt")
	if err != nil {
		log.Panic(err)
	}

	bot := app.NewPollBot(api, pollOptsProv)
	bot.Run(updateConfig)
}
