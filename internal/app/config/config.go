package config

// Интерфейс работы с конфигурацией бота
type Config interface {
	// Добавить вариант опроса
	AddPollOpt(opt string) error
	// Удалить вариант опроса
	DeletePollOpt(opt string) error
	// Вернуть список вариантов опроса
	GetPollOpts() ([]string, error)

	// Получить время жизни опроса
	GetPollLifetime() (int, error)
	// Установить время жизни опроса
	SetPollLifetime(timeout int) error

	// Добавить чат, в котором бот будет создавать опрос
	AddChat(chat int64) error
	// Удалить чат
	DeleteChat(chat int64) error
	// Вернуть список чатов, в которых будет создаваться опрос
	GetChats() ([]int64, error)

	// Установка паттерна для планировщика
	SetCronPattern(ptrn string) error
	// Получить паттерн планировщика
	GetCronPattern() (string, error)
}
