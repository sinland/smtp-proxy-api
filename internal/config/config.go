package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Port      int    `yaml:"port"`
		JwtSecret string `yaml:"jwt_secret"`
		ApiKey    string `yaml:"api_key"`
		BotToken  string `yaml:"bot_token"`
	} `yaml:"server"`
	SMTP struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
}

func New(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
