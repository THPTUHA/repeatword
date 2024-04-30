package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Configs struct {
	DB struct {
		Mysql struct {
			Username     string `yaml:"username"`
			Password     string `yaml:"password"`
			Port         int    `yaml:"port"`
			URI          string `yaml:"uri"`
			DatabaseName string `yaml:"databaseName"`
		}
	}
	DataDir string `yaml:"dataDir"`
}

var config *Configs

func Get() (*Configs, error) {
	if config == nil {
		return nil, errors.New("empty config")
	}
	return config, nil
}

func Set(f string) (*Configs, error) {
	config = &Configs{}
	file, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	err = GetYaml(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func GetYaml(f []byte, s interface{}) error {
	y := yaml.Unmarshal(f, s)
	return y
}
