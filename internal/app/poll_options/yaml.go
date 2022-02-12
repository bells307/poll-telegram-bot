package poll_options

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type YamlPollOptionsProvider struct {
	path string
}

func NewYamlPollOptionsProvider(path string) (*YamlPollOptionsProvider, error) {
	prov := YamlPollOptionsProvider{path}

	// Создаем файл, если его нет
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return &prov, nil
}

func (p YamlPollOptionsProvider) Add(opt string) error {
	list, err := p.List()
	if err != nil {
		return nil
	}

	list = append(list, opt)

	data, err := yaml.Marshal(list)
	if err != nil {
		return err
	}

	err = os.WriteFile(p.path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p YamlPollOptionsProvider) Delete(opt string) error {
	list, err := p.List()
	if err != nil {
		return nil
	}

	found_idx := -1
	for idx, val := range list {
		if val == opt {
			found_idx = idx
			break
		}
	}

	if found_idx == -1 {
		return errors.New("option not found")
	}

	// Удаляем элемент из списка
	list = append(list[:found_idx], list[found_idx+1:]...)

	data, err := yaml.Marshal(list)
	if err != nil {
		return err
	}

	err = os.WriteFile(p.path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p YamlPollOptionsProvider) List() ([]string, error) {
	data, err := ioutil.ReadFile(p.path)
	if err != nil {
		return nil, err
	}

	var list []string
	err = yaml.Unmarshal(data, list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
