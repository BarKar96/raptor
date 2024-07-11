package config

import (
	"context"
	"sync"
	"time"

	"github.com/barkar96/raptor/libs/logging"

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

	RedisAddr     []string      `split_words:"true" default:"127.0.0.1:5433"`
	RedisPassword string        `split_words:"true" default:""`
	RedisTimeout  time.Duration `split_words:"true" default:"3s"`
}
