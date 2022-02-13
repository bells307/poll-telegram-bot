package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type (
	// Структура работы с yaml-конфигурацией бота
	YamlConfigProvider struct {
		*sync.Mutex
		// Путь к файлу конфигурации
		path string
	}

	// yaml-конфигурация
	ConfigDesc struct {
		// Список чатов, в которых состоит бот
		Chats []string `yaml:"chats"`
		// Список вариантов голосования
		PollOptions []string `yaml:"poll_options"`
	}
)

func NewYamlConfigProvider(path string) (*YamlConfigProvider, error) {
	prov := YamlConfigProvider{path: path, Mutex: &sync.Mutex{}}

	// Проверяем существование файла конфигурации
	if _, err := os.Stat(path); err == nil {
		// Файл существует, проверяем его валидность
		log.Printf("Found an existing configuration file at %s\n", path)

		_, err := prov.unmarshalConfig()
		if err != nil {
			// Текущий файл конфигурации невалидный, создаем конфигурацию по умолчанию
			log.Printf("The existing configuration is not valid: %v. Creating new empty configuration\n", err)
			if err := prov.createDefaultConfig(); err != nil {
				return nil, err
			}
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// Файл не существует, создаем конфигурацию по умолчанию
		log.Printf("The configuration file %s is not found. Creating new empty configuration\n", path)

		err = prov.createDefaultConfig()
		if err != nil {
			return nil, err
		}
	} else {
		// Файл Шредингера: возможно существует, а возможно и нет. Выведем в лог и вернем ошибку
		log.Printf("Configuration file %s processing error: %v\n", path, err)
		return nil, err
	}

	return &prov, nil
}

// Добавить вариант опроса
func (p YamlConfigProvider) AddPollOpt(opt string) error {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	cfg, err := p.unmarshalConfig()
	if err != nil {
		return err
	}

	// Проверяем уникальность
	for _, val := range cfg.PollOptions {
		if val == opt {
			return fmt.Errorf("option \"%s\" already exists", opt)
		}
	}

	cfg.PollOptions = append(cfg.PollOptions, opt)

	return p.marshalConfig(*cfg)
}

// Удалить вариант опроса
func (p YamlConfigProvider) DeletePollOpt(opt string) error {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	cfg, err := p.unmarshalConfig()
	if err != nil {
		return err
	}

	foundIdx := -1
	for idx, val := range cfg.PollOptions {
		if val == opt {
			foundIdx = idx
			break
		}
	}

	if foundIdx == -1 {
		return errors.New("option not found")
	}

	// Удаляем элемент из списка
	cfg.PollOptions = append(cfg.PollOptions[:foundIdx], cfg.PollOptions[foundIdx+1:]...)

	return p.marshalConfig(*cfg)
}

// Вернуть список вариантов опроса
func (p YamlConfigProvider) ListPollOpt() ([]string, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	cfg, err := p.unmarshalConfig()
	if err != nil {
		return nil, err
	}

	return cfg.PollOptions, nil
}

// Получить объект конфигурации из файла
func (p *YamlConfigProvider) unmarshalConfig() (*ConfigDesc, error) {
	data, err := ioutil.ReadFile(p.path)
	if err != nil {
		return nil, err
	}

	var cfg ConfigDesc
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Записать объект конфигурации в файл
func (p *YamlConfigProvider) marshalConfig(cfg ConfigDesc) error {
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(p.path, data, 0644)
}

// Создать конфигурацию по умолчанию
func (p *YamlConfigProvider) createDefaultConfig() error {
	cfg := ConfigDesc{
		Chats:       []string{},
		PollOptions: []string{},
	}

	if err := p.marshalConfig(cfg); err != nil {
		log.Printf("Error while creating default configuration file: %v", err)
		return err
	}
	return nil
}
