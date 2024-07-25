package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	fiberRecover "github.com/gofiber/fiber/v3/middleware/recover"
	"net/http"
	"trading-ace/model"
)

func panicRecover(config ...fiberRecover.Config) fiber.Handler {
	cfg := fiberRecover.ConfigDefault

	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c fiber.Ctx) (retErr error) {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				var appErr *model.AppError
				if errors.As(err, &appErr) {
					c.Status(appErr.HTTPStatusCode)
					_ = c.JSON(appErr)
					return
				}

				retErr = c.Status(http.StatusInternalServerError).JSON(model.NewServerError())
			}
		}()

		return c.Next()
	}
}
