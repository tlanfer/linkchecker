package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const configFile = "config.yaml"

type Config struct {
	LogFile  string        `yaml:"log_file"`
	Interval time.Duration `yaml:"interval"`
	Hosts    []string      `yaml:"hosts"`
}

func Load() Config {
	f, err := os.OpenFile(configFile, os.O_RDONLY, os.ModePerm)
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) {
			defaultConfig := Config{
				"log.csv",
				15 * time.Second,
				[]string{"youtube.com", "twitch.com"},
			}
			Save(defaultConfig)
			return defaultConfig
		}
	}

	c := Config{}
	yaml.NewDecoder(f).Decode(&c)
	return c
}

func Save(config Config) {
	f, err := os.OpenFile(configFile, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	yaml.NewEncoder(f).Encode(&config)
}
