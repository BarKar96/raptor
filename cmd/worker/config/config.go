package config

import (
	"context"
	"sync"

	"github.com/barkar96/worker/libs/logging"

	"github.com/kelseyhightower/envconfig"
)

var (
	instance *Config
	once     sync.Once
)

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
		err := envconfig.Process("", instance)
		if err != nil {
			logging.WithFatalError(context.TODO(), err, "failed to read configuration")
		}
	})
	return instance
}

type Config struct {
	ApplicationName string `split_words:"true" default:"worker"`
	EnvironmentName string `split_words:"true" default:"local"`
}
