package raptor

import (
	"github.com/barkar96/raptor/lib/api"
	"github.com/barkar96/raptor/lib/logging"
	"github.com/gofiber/fiber/v2"
)

var _ api.API = (*UserAPI)(nil)

type UserAPI struct{}

func NewUserAPI() UserAPI {
	return UserAPI{}
}

func (u UserAPI) Register(app *fiber.App) error {
	apiV1 := app.Group("/api/v1")

	apiV1.Post("/user", u.HandleCreateUser)

	return nil
}

func (u UserAPI) HandleCreateUser(c *fiber.Ctx) error {
	ctx := c.Context()
	logging.Info(ctx, "hello world")
	return nil
}
