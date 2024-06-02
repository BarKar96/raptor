package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/barkar96/worker/cmd"
	"github.com/barkar96/worker/libs/actor"
)

func main() {
	cmd.SetGoMaxProcs()
	ctx := context.Background()
	cfg := map[string]string{}

	if err := run(ctx, cfg); err != nil {
		slog.Error("TODO")
	}
}

func run(ctx context.Context, sth interface{}) error {
	g := actor.New()
	g.Run(ctx, 5*time.Second)

	return nil
}
