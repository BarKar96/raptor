package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/barkar96/raptor/cmd"
	"github.com/barkar96/raptor/cmd/worker/config"
	"github.com/barkar96/raptor/lib/actor"
	"github.com/barkar96/raptor/lib/api"
	"github.com/barkar96/raptor/lib/logging"
	"github.com/barkar96/raptor/lib/redis"

	"github.com/barkar96/raptor/api/raptor"
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
	logging.Info(ctx, "Redis initialized")
	g.Add(rdb)

	// === PostgreSQL ===

	// === Fiber ===
	userAPI := raptor.NewUserAPI()
	app, err := api.NewAPI("0.0.0.0:8080", "dev", &userAPI)
	if err != nil {
		logging.WithFatalError(ctx, err, "API initilization failed")
	}
	g.Add(app)

	err = g.Run(ctx, 5*time.Second)
	if err != nil {
		logging.WithError(ctx, err, "application stopped on error")
	}

	return nil
}
