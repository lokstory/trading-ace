package api

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v3"
	"log"
	"trading-ace/model"
)

func (a *App) ListUserPointTaskStates(c fiber.Ctx) error {
	address := c.Params("address")

	if !common.IsHexAddress(address) {
		return model.NewParameterError()
	}

	addr := common.HexToAddress(address)

	log.Println("list point task address:", addr)

	ret, appErr := a.userPointTaskStateRepository.ListWithTaskByAddress(c.Context(), addr.Hex())
	if appErr != nil {
		return appErr
	}

	log.Println("list point task ret:", ret)

	return successResponse(c, ret)
}
