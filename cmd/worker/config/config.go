package config

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/barkar96/worker/libs/logging"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ApplicationName string `split_words:"true" default:"worker"`
	EnvironmentName string `split_words:"true" default:"local"`
}

func ParseConfig(srvName string) *Config {
	var cfg Config
	err := envconfig.Process(srvName, &cfg)
	if err != nil {
		logging.WithFatalError(context.TODO(), err, "failed to read configuration")
	}
	return &cfg
}

var (
	instance *Config
	once     sync.Once
)

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
		err := envconfig.Process("", instance)
		if err != nil {
			slog.Error("read configuration error", err)
			os.Exit(1)
		}
	})
	return instance
}
