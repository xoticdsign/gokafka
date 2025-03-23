package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Client Client `yaml:"client"`

	KafkaConfig KafkaConfig `yaml:"kafka"`
}

type Client struct {
	Timeout time.Duration `yaml:"timeout"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
	Topics  []string `yaml:"topics"`
}

func LoadConfig(env string) (Config, error) {
	var cfg Config

	switch env {
	case "prod":

	case "dev":

	default:
		err := cleanenv.ReadConfig("./config/notificator/local.yaml", &cfg)
		if err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}
