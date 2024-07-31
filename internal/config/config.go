package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

const APP_ENVIRONMENT_LOCAL string = "LOCAL"
const APP_ENVIRONMENT_TEST string = "TEST"
const APP_ENVIRONMENT_PRODUCTION string = "PRODUCTION"

var defaultConfig = configuration{}

type configuration struct {
	Discord     discordConfig `yaml:"discord"`
	Mumble      mumbleConfig  `yaml:"mumble"`
	Prefix      string        `yaml:"prefix"`
	Status      string        `yaml:"status"`
	Environment string        `yaml:"environment"`
}

type discordConfig struct {
	Token  string `yaml:"token"`
	Enable bool   `yaml:"enable"`
}

type mumbleConfig struct {
	Enable   bool   `yaml:"enable"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	TLS      bool   `yaml:"tls"`
	Username string `yaml:"username"`
}

var config *configuration

func Load() {
	file, err := os.Open("config.yml")
	if err != nil {
		slog.Error("failed to open config file", "error", err)
		return
	}
	defer file.Close()

	config = &defaultConfig
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		slog.Error("failed to decode config file", "error", err)
		return
	}
}

func GetEnvironment() string {
	return config.Environment
}

func IsEnvironment(environments ...string) bool {
	if len(environments) == 0 {
		return config.Environment == environments[0]
	}

	for _, environment := range environments {
		if config.Environment == environment {
			return true
		}
	}

	return false
}

func GetBotPrefix() string {
	return config.Prefix
}

func GetBotStatus() string {
	return config.Status
}

func IsDiscordEnabled() bool {
	return config.Discord.Enable
}

func GetDiscordToken() string {
	return config.Discord.Token
}

func IsMumbleEnabled() bool {
	return config.Mumble.Enable
}

func GetMumbleHost() string {
	return config.Mumble.Host
}

func GetMumblePort() int {
	return config.Mumble.Port
}

func GetMumbleUsername() string {
	return config.Mumble.Username
}

func GetMumbleTLS() bool {
	return config.Mumble.TLS
}
