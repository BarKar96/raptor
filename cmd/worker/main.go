package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/barkar96/worker/cmd"
	"github.com/barkar96/worker/cmd/worker/config"
	"github.com/barkar96/worker/libs/actor"
	"github.com/barkar96/worker/libs/logging"
	"github.com/barkar96/worker/libs/redis"
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
	_, err := redis.New([]string{"127.0.0.1:5433"}, "", time.Second*5)
	if err != nil {
		logging.WithFatalError(ctx, err, "Redis initilization failed")
	}

	// === PostgreSQL ===

	logging.Info(ctx, cfg.ApplicationName)

	g.Run(ctx, 5*time.Second)

	return nil
}
