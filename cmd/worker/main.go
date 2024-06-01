package main

import (
	"context"
	"log/slog"

	"github.com/barkar96/worker/cmd"
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

	return nil
}
