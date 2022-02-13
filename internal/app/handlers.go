package app

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Помощь
func (b *PollBot) HelpHandler(u *tm.Update) {
	b.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, `
Bot commands:
/poll - create a poll

/add [options] - add poll options (separated by ;)
/del [options] - delete poll options (separated by ;)
/list - get a list of poll options

/lifetime <value> - set poll lifetime (secs)
`))
}

// Создание опроса
func (b *PollBot) PollHandler(u *tm.Update) {
	chats, err := b.Config.ListChats()
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Can't get bot chats while creating poll: %v", err))
		return
	}

	opts, err := b.Config.ListPollOpt()
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Error getting poll options list: %v", err))
		return
	}

	lifetime, err := b.Config.GetPollLifetime()
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Error getting poll lifetime: %v", err))
		return
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
			b.AnswerAndLog(u, fmt.Sprintf("Error creating poll for chat %v: %v", chat, err))
		}
	}

	// b.AnswerAndLog(u, "Polls created")
}

// Добавление варианта опроса
func (b *PollBot) AddPollOptionHandler(u *tm.Update) {
	args, err := getCommandArgs(u)
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("%v", err))
		return
	}

	opts := strings.Split(strings.Join(args, " "), "; ")

	for _, opt := range opts {
		err := b.Config.AddPollOpt(opt)

		if err != nil {
			b.AnswerAndLog(u, fmt.Sprintf("Error adding poll option \"%s\": %v", opt, err))
		} else {
			b.AnswerAndLog(u, fmt.Sprintf("The poll option \"%s\" added", opt))
		}
	}
}

// Удаление варианта опроса
func (b *PollBot) DeletePollOptionHandler(u *tm.Update) {
	args, err := getCommandArgs(u)
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("%v", err))
		return
	}

	opts := strings.Split(strings.Join(args, " "), "; ")

	for _, opt := range opts {
		err := b.Config.DeletePollOpt(opt)

		if err != nil {
			b.AnswerAndLog(u, fmt.Sprintf("Error deleting poll option \"%s\": %v", opt, err))
		} else {
			b.AnswerAndLog(u, fmt.Sprintf("The poll option \"%s\" deleted", opt))
		}
	}
}

// Получить список вариантов опроса
func (b *PollBot) ListPollOptionsHandler(u *tm.Update) {
	list, err := b.Config.ListPollOpt()
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Error getting poll options list: %v", err))
		return
	}

	b.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID,
		fmt.Sprintf("Current list of poll options: %s", strings.Join(list, ", ")),
	))
}

// Установка времени жизни опросов
func (b *PollBot) LifetimeHandler(u *tm.Update) {
	args, err := getCommandArgs(u)
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("%v", err))
		return
	}

	val, err := strconv.Atoi(args[0])
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Can't set poll lifetime: %v", err))
		return
	}

	if err = b.Config.SetPollLifetime(val); err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Can't set poll lifetime: %v", err))
		return
	}

	b.AnswerAndLog(u, fmt.Sprintf("New lifetime value is set: %v seconds", val))
}

// Запись сообщения в лог и ответ запросившему
func (b *PollBot) AnswerAndLog(u *tm.Update, msg string) {
	log.Println(msg)
	b.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg))
}

// Получение аргументов, переданных в команду
func getCommandArgs(u *tm.Update) ([]string, error) {
	args, ok := u.Context["args"]
	if !ok {
		return nil, errors.New("can't get args from update")
	}

	var argsSlice []string
	switch reflect.TypeOf(args).Kind() {
	case reflect.Slice:
		v := reflect.ValueOf(args)
		argsSlice = v.Interface().([]string)
	}

	return argsSlice, nil
}
