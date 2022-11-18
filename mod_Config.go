package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	SMPP struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"smpp"`

	SMTP struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Encr int    `yaml:"encr"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		From string `yaml:"from"`
	} `yaml:"smtp"`

	Telegram struct {
		BotToken string `yaml:"botToken"`
	} `yaml:"telegram"`
}

func readConfigFile(cfg *Config) {
	f, err := os.Open("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
