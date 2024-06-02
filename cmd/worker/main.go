package main

import (
	"context"
	"time"

	"github.com/barkar96/worker/cmd"
	"github.com/barkar96/worker/libs/actor"
	"github.com/barkar96/worker/libs/logging"
)

func main() {
	cmd.SetGoMaxProcs()
	ctx := context.Background()
	cfg := map[string]string{}

	if err := run(ctx, cfg); err != nil {
		logging.WithFatalError(ctx, err, "application stopped on error")
	}
}

func run(ctx context.Context, sth interface{}) error {
	g := actor.New()
	g.Run(ctx, 5*time.Second)

	return nil
}
