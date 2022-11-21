package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const configFile = "config.yaml"

type Config struct {
	Prefix    string        `yaml:"prefix"`
	Interval  time.Duration `yaml:"interval"`
	Threshold Threshold     `json:"threshold"`
	Hosts     []string      `yaml:"hosts"`
}

type Threshold struct {
	PacketLoss float64       `yaml:"packet_loss"`
	Rtt        time.Duration `json:"rtt"`
}

func Load() Config {
	c := Config{
		"monitor_",
		15 * time.Second,
		Threshold{
			PacketLoss: 0,
			Rtt:        100 * time.Millisecond,
		},
		[]string{"youtube.com", "twitch.tv", "regjeringen.no"},
	}

	f, err := os.OpenFile(configFile, os.O_RDONLY, os.ModePerm)
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) {
			Save(c)
			return c
		}
	}

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
