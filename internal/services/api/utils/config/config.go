package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	Client       Client        `yaml:"client"`

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
		// TODO

	case "dev":
		// TODO

	default:
		err := cleanenv.ReadConfig("./config/api/local.yaml", &cfg)
		if err != nil {
			return Config{}, err
		}

	}

	return cfg, nil
}
