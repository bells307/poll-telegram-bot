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
/cron - set cron pattern
`))
}

// Создание опроса
func (b *PollBot) PollHandler(u *tm.Update) {
	if err := b.createPolls(); err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Error creating polls: %v", err))
	}
}

// Добавление варианта опроса
func (b *PollBot) AddPollOptionHandler(u *tm.Update) {
	args, err := getCommandArgs(u)
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("%v", err))
		return
	}

	opts := strings.Split(strings.Join(args, " "), ";")

	for _, opt := range opts {
		err := b.Config.AddPollOpt(strings.TrimSpace(opt))

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

	opts := strings.Split(strings.Join(args, " "), ";")

	for _, opt := range opts {
		err := b.Config.DeletePollOpt(strings.TrimSpace(opt))

		if err != nil {
			b.AnswerAndLog(u, fmt.Sprintf("Error deleting poll option \"%s\": %v", opt, err))
		} else {
			b.AnswerAndLog(u, fmt.Sprintf("The poll option \"%s\" deleted", opt))
		}
	}
}

// Получить список вариантов опроса
func (b *PollBot) ListPollOptionsHandler(u *tm.Update) {
	list, err := b.Config.GetPollOpts()
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Error getting poll options list: %v", err))
		return
	}

	var res string
	if len(list) == 0 {
		res = "The list of poll options is empty"
	} else {
		res += "Current list of poll options:\n"
		for _, opt := range list {
			res += " - " + opt + "\n"
		}
	}

	b.BotAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, res))
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

func (b *PollBot) CronHandler(u *tm.Update) {
	args, err := getCommandArgs(u)
	if err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("%v", err))
		return
	}

	pat := strings.Join(args, " ")
	if err := b.Config.SetCronPattern(pat); err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Error setting cron pattern: %v", err))
	}

	if err := b.startScheduler(); err != nil {
		b.AnswerAndLog(u, fmt.Sprintf("Error updating scheduler: %v", err))
	} else {
		b.AnswerAndLog(u, "Cron scheduler updated")
	}
}

// Обработка добавления бота в чат
func (b *PollBot) SelfAddedToChat(u *tm.Update) {
	log.Printf("Bot added to chat %v %v", u.EffectiveChat().ID, u.EffectiveChat().Title)
	b.Config.AddChat(u.EffectiveChat().ID)
}

// Обработка удаления бота из чата
func (b *PollBot) SelfRemovedFromChat(u *tm.Update) {
	log.Printf("Bot removed from chat %v %v", u.EffectiveChat().ID, u.EffectiveChat().Title)
	b.Config.DeleteChat(u.EffectiveChat().ID)
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
