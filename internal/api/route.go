package api

import (
	"github.com/gofiber/fiber/v3"
)

func (a *App) initRoute(app *fiber.App) {
	userGroup := app.Group("/users")

	userGroup.Get("/:address/point-task-states", a.ListUserPointTaskStates)
	userGroup.Get("/:address/point-histories", a.ListUserPointHistories)
}
