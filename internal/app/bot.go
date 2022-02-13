package app

import (
	tm "github.com/and3rson/telemux/v2"
	"github.com/bells307/poll-telegram-bot/internal/app/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PollBot struct {
	BotAPI *tgbotapi.BotAPI
	Config config.Config
}

func NewPollBot(api *tgbotapi.BotAPI, config config.Config) *PollBot {
	bot := PollBot{BotAPI: api, Config: config}
	return &bot
}

// Запуск обработки апдейтов, приходящих к боту
func (b *PollBot) Run(updConfig tgbotapi.UpdateConfig) {
	updChan := b.BotAPI.GetUpdatesChan(updConfig)
	mux := b.newMux()

	for upd := range updChan {
		mux.Dispatch(b.BotAPI, upd)
	}
}

// Маршрутизатор функций обработки команд
func (b *PollBot) newMux() *tm.Mux {
	mux := tm.NewMux().
		AddHandler(tm.NewCommandHandler(
			"help",
			nil,
			func(u *tm.Update) {
				b.HelpHandler(u)
			},
		)).
		AddHandler(tm.NewCommandHandler(
			"poll",
			nil,
			func(u *tm.Update) {
				b.PollHandler(u)
			},
		)).
		AddHandler(tm.NewCommandHandler(
			"add",
			nil,
			func(u *tm.Update) {
				b.AddPollOptionHandler(u)
			},
		)).
		AddHandler(tm.NewCommandHandler(
			"del",
			nil,
			func(u *tm.Update) {
				b.DeletePollOptionHandler(u)
			},
		)).
		AddHandler(tm.NewCommandHandler(
			"list",
			nil,
			func(u *tm.Update) {
				b.ListPollOptionsHandler(u)
			},
		))
		// AddHandler(tm.NewHandler(
		// 	nil,
		// 	func(u *tm.Update) {
		// 		b.Api.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "You said: "+u.Message.Text))
		// 	},
		// ))
	return mux
}
