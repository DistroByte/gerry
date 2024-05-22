package bot

import (
	"fmt"
	"os"

	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/yaml.v3"
)

var defaultConfig = Config{}

type Config struct {
	Discord DiscordConfig `yaml:"discord"`
	Mumble  MumbleConfig  `yaml:"mumble"`
	Prefix  string        `yaml:"prefix"`
	Markov  MarkovConfig  `yaml:"markov"`
}

type DiscordConfig struct {
	Token    string         `yaml:"token"`
	GuildIDs []snowflake.ID `yaml:"guild_ids"`
	Username string         `yaml:"username"`
}

type MumbleConfig struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	TLS      bool   `yaml:"tls"`
	Username string `yaml:"username"`
}

type MarkovConfig struct {
	FlushInterval int `yaml:"flush_interval"`
	Gravity       int `yaml:"gravity"`
}

func ReadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	cfg := defaultConfig
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to decode config: %w", err)
	}

	return cfg, nil
}
