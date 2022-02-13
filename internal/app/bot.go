package app

import (
	"log"

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
			"start",
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
		)).
		AddHandler(tm.NewCommandHandler(
			"lifetime",
			nil,
			func(u *tm.Update) {
				b.LifetimeHandler(u)
			},
		)).
		AddHandler(tm.NewHandler(
			func(u *tm.Update) bool {
				if msg := u.EffectiveMessage(); msg != nil {
					if members := msg.NewChatMembers; members != nil {
						for _, m := range members {
							if m.ID == b.BotAPI.Self.ID {
								return true
							}
						}
					}
				}

				return false
			},
			func(u *tm.Update) {
				b.SelfAddedToChat(u)
			},
		)).
		AddHandler(tm.NewHandler(
			func(u *tm.Update) bool {
				if msg := u.EffectiveMessage(); msg != nil {
					if member := msg.LeftChatMember; member != nil {
						if member.ID == b.BotAPI.Self.ID {
							return true
						}
					}
				}

				return false
			},
			func(u *tm.Update) {
				b.SelfRemovedFromChat(u)
			},
		))
	return mux
}

func (b *PollBot) createPolls() error {
	chats, err := b.Config.GetChats()
	if err != nil {
		return err
	}

	opts, err := b.Config.GetPollOpts()
	if err != nil {
		return err
	}

	lifetime, err := b.Config.GetPollLifetime()
	if err != nil {
		return err
	}

	for _, chat := range chats {
		pollConfig := tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: chat,
			},
			Question:    "Question",
			Options:     opts,
			IsAnonymous: false,
			OpenPeriod:  lifetime,
		}

		_, err := b.BotAPI.Send(pollConfig)
		if err != nil {
			log.Printf("Error creating poll for chat %v: %v", chat, err)
		}
	}

	return nil
}
