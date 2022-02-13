package config

// Интерфейс работы с конфигурацией бота
type Config interface {
	// Добавить вариант опроса
	AddPollOpt(opt string) error
	// Удалить вариант опроса
	DeletePollOpt(opt string) error
	// Вернуть список вариантов опроса
	ListPollOpt() ([]string, error)
}
