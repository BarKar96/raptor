package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/barkar96/raptor/cmd"
	"github.com/barkar96/raptor/cmd/raptor/config"
	"github.com/barkar96/raptor/libs/actor"
	"github.com/barkar96/raptor/libs/logging"
	"github.com/barkar96/raptor/libs/redis"
)

func main() {
	cmd.SetGoMaxProcs()
	ctx := context.Background()
	cfg := config.GetInstance()
	logging.Init(slog.LevelInfo, cfg.EnvironmentName, cfg.ApplicationName, false)

	if err := run(ctx, cfg); err != nil {
		logging.WithFatalError(ctx, err, "application stopped on error")
	}
}

func run(ctx context.Context, cfg *config.Config) error {
	g := actor.New()

	// === Redis ===
	rdb, err := redis.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisTimeout)
	if err != nil {
		logging.WithFatalError(ctx, err, "Redis initilization failed")
	}
	g.Add(rdb)

	// === PostgreSQL ===

	logging.Info(ctx, cfg.ApplicationName)

	g.Run(ctx, 5*time.Second)

	return nil
}
