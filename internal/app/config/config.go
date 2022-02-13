package config

// Интерфейс работы с конфигурацией бота
type Config interface {
	// Добавить вариант опроса
	AddPollOpt(opt string) error
	// Удалить вариант опроса
	DeletePollOpt(opt string) error
	// Вернуть список вариантов опроса
	ListPollOpt() ([]string, error)
	// Получить время жизни опроса
	GetPollLifetime() (int, error)
	// Установить время жизни опроса
	SetPollLifetime(timeout int) error
	// Вернуть список чатов, в которых будет создаваться опрос
	ListChats() ([]int64, error)
}
