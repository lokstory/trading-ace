package api

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v3"
	"trading-ace/model"
)

func (a *App) ListUserPointHistories(c fiber.Ctx) error {
	address := c.Params("address")

	if !common.IsHexAddress(address) {
		return model.NewParameterError()
	}

	addr := common.HexToAddress(address)

	ret, appErr := a.userPointHistoryRepository.ListWithTaskByAddress(c.Context(), addr.Hex())
	if appErr != nil {
		return appErr
	}

	return successResponse(c, ret)
}
