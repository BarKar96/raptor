package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/barkar96/worker/cmd"
	"github.com/barkar96/worker/cmd/worker/config"
	"github.com/barkar96/worker/libs/actor"
	"github.com/barkar96/worker/libs/logging"
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

func run(ctx context.Context, config *config.Config) error {
	g := actor.New()

	logging.Info(ctx, "hello world")

	g.Run(ctx, 5*time.Second)

	return nil
}
