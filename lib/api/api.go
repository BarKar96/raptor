package api

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/barkar96/raptor/lib/actor"
	"github.com/barkar96/raptor/lib/logging"

	"github.com/gofiber/fiber/v2"
)

var _ actor.Actor = (*BaseAPI)(nil)

type BaseAPI struct {
	App  *fiber.App
	addr string
}

type API interface {
	Register(app *fiber.App) error
}

func NewAPI(addr, env string, apiList ...API) (*BaseAPI, error) {
	cfg := fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Immutable:    true,
	}

	ba := &BaseAPI{
		App:  fiber.New(cfg),
		addr: addr,
	}

	for _, apiEntry := range apiList {
		err := apiEntry.Register(ba.App)
		if err != nil {
			return nil, fmt.Errorf("register API: %w", err)
		}
	}

	return ba, nil
}

func (b *BaseAPI) Name() string {
	return "API"
}

func (b *BaseAPI) Start(_ context.Context) error {
	if strings.HasPrefix(b.addr, "unix:") {
		return b.listenAndServeUNIX()
	}
	return b.ListenAndServe()
}

func (b *BaseAPI) Stop(ctx context.Context) {
	err := b.App.ShutdownWithContext(ctx)
	if err != nil {
		logging.WithError(ctx, err, "shutdown timeout exceeded")
	}
}

func (b *BaseAPI) ListenAndServe() error {
	return b.App.Listen(b.addr)
}

func (b *BaseAPI) listenAndServeUNIX() error {
	addr := b.addr[5:]
	mode := os.ModeSocket | os.ModePerm
	ln, err := net.Listen("unix", addr)
	if err != nil {
		return err
	}

	if err = os.Chmod(addr, mode); err != nil {
		return fmt.Errorf("cannot chmod %#o for %q: %w", mode, addr, err)
	}

	return b.App.Listener(ln)
}
