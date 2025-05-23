package config

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const APP_ENVIRONMENT_LOCAL string = "LOCAL"
const APP_ENVIRONMENT_TEST string = "TEST"
const APP_ENVIRONMENT_PRODUCTION string = "PROD"

var defaultConfig = configuration{}

var ShutdownChannel = make(chan os.Signal, 1)
var StartTime time.Time

type configuration struct {
	Discord     discordConfig `yaml:"discord"`
	Mumble      mumbleConfig  `yaml:"mumble"`
	HTTP        httpConfig    `yaml:"http"`
	Prefix      string        `yaml:"prefix" default:">"`
	Status      string        `yaml:"status"`
	Environment string        `yaml:"environment" default:"LOCAL" validate:"required,oneof=LOCAL TEST PROD"`
	Domain      string        `yaml:"domain"`
	Name        string        `yaml:"name"`
}

type discordConfig struct {
	Token  string `yaml:"token"`
	Enable bool   `yaml:"enable" default:"false"`
}

type httpConfig struct {
	Port   int  `yaml:"port" default:"8080"`
	Enable bool `yaml:"enable" default:"false"`
}

type mumbleConfig struct {
	Enable   bool   `yaml:"enable" default:"false"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port" default:"64738"`
	TLS      bool   `yaml:"tls" default:"false"`
	Username string `yaml:"username"`
}

var config *configuration

func Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Msg("failed to open config file, creating empty file")
		Generate(path)
		return err
	}
	defer file.Close()

	config = &defaultConfig
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		log.Error().Err(err).Msg("failed to decode config file")
		return err
	}

	log.Info().Str("file", path).Msg("config file loaded successfully")
	return nil
}

func Generate(filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		log.Error().Err(err).Msg("failed to create config file")
		return
	}
	defer file.Close()

	log.Info().Str("file", filepath).Msg("writing default config to file")
	if err := yaml.NewEncoder(file).Encode(defaultConfig); err != nil {
		log.Error().Err(err).Msg("failed to encode config file")
		return
	}

	log.Info().Msg("config file generated successfully")
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

func GetDomain() string {
	return config.Domain
}

func GetBotPrefix() string {
	return config.Prefix
}

func GetBotStatus() string {
	return config.Status
}

func GetBotName() string {
	return config.Name
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

func IsHTTPEndpointEnabled() bool {
	return config.HTTP.Enable
}

func GetHTTPPort() int {
	return config.HTTP.Port
}
