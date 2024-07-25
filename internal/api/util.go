package api

import (
	"github.com/gofiber/fiber/v3"
	"trading-ace/model"
)

func successResponse(c fiber.Ctx, body interface{}) error {
	return c.JSON(model.NewOKResponse(body))
}
