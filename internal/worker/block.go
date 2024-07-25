package worker

import (
	"fmt"
	"math/big"
	"time"
	"trading-ace/model"
)

func (a *App) getOrUpdateBlockTime(blockNumber uint64) (time.Time, *model.AppError) {
	blockTime, apiErr := a.blockRepository.GetBlockTime(a.ctx, blockNumber)

	if apiErr != nil {
		blockNumberValue := new(big.Int).SetUint64(blockNumber)
		block, err := a.rpcClient.BlockByNumber(a.ctx, blockNumberValue)
		if err != nil {
			return blockTime, model.NewAppError(model.OperateFailed).Err(fmt.Errorf("eth get block by number error: %w", err))
		}

		blockTime = time.Unix(int64(block.Time()), 0)
		_ = a.blockRepository.Create(a.ctx, blockNumber, blockTime)
	}

	return blockTime, nil
}
