package actor

import (
	"context"
	"log/slog"
	"time"
)

type Actor interface {
	Start(context.Context) error
	Stop(context.Context)
	Name() string
}

type Group struct {
	actors []Actor
}

func New() *Group {
	g := &Group{}
	// g.Add(signalCatcher{})
	return g
}

func (g *Group) Add(a Actor) {
	g.actors = append(g.actors, a)
}

func (g *Group) Run(ctx context.Context, shutDownTimeout time.Duration) error {
	if len(g.actors) == 0 {
		return nil
	}

	slog.Info("application starting")
	startCtx, cancel := context.WithCancel(ctx)
	errors := make(chan error, len(g.actors))
	for _, a := range g.actors {
		go func(a Actor) {
			slog.Info("actor starting")
			errors <- a.Start(startCtx)
		}(a)
	}

	// wait for first error
	err := <-errors

	slog.Info("initiating shutdown", slog.Duration("timeout", shutDownTimeout))
	cancel()

	stopCtx, cancel := context.WithTimeout(ctx, shutDownTimeout)
	defer cancel()

	// stop all goroutines
	for i := len(g.actors) - 1; i >= 0; i-- {
		go func(i int) {
			a := g.actors[i]
			a.Stop(stopCtx)
			slog.Info("stopped")
		}(i)
	}

	// listen for shutdown timeout
	go func() {
		<-stopCtx.Done()
		if stopCtx.Err() == context.DeadlineExceeded {
			slog.Error("shutdown timeout exceeded")
		}
	}()

	// wait for all goroutines to stop
	for i := 1; i < cap(errors); i++ {
		<-errors
	}

	return err

}
