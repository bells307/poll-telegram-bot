package main

import (
	"flag"
	"log"

	"github.com/bells307/poll-telegram-bot/internal/app"
	"github.com/bells307/poll-telegram-bot/internal/app/poll_options"
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
	api, err := tgbotapi.NewBotAPI("TOKEN")
	if err != nil {
		log.Panic(err)
	}

	api.Debug = true

	log.Printf("Authorized on account %s", api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	var pollOptsProv poll_options.PollOptionsProvider = nil
	if mode == "yaml" {
		pollOptsProv, err = poll_options.NewYamlPollOptionsProvider("opts.yaml")
		if err != nil {
			log.Panic(err)
		}
	}

	bot := app.NewPollBot(api, pollOptsProv)
	bot.Run(updateConfig)
}
