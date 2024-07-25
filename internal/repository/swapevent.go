package repository

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"trading-ace/model"
)

type SwapEventRepository struct {
	*BaseRepository
}

func (r *SwapEventRepository) SaveSwapEvents(ctx context.Context, contractABI abi.ABI, blockTime time.Time, logs []types.Log) *model.AppError {
	return r.transaction(ctx, func(tx pgx.Tx) *model.AppError {
		for _, item := range logs {
			var swapLog model.SwapLog

			err := contractABI.UnpackIntoInterface(&swapLog, "Swap", item.Data)

			swapLog.Sender = common.HexToAddress(item.Topics[1].Hex())
			swapLog.To = common.HexToAddress(item.Topics[2].Hex())

			_, err = tx.Exec(ctx, `
INSERT INTO public.swap_event (contract_address, block_number, log_index, sender, to_address,
                        amount0_in, amount1_in, amount0_out, amount1_out, block_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (contract_address, block_number, log_index) DO NOTHING
	`, item.Address.Hex(), item.BlockNumber, item.Index, swapLog.Sender.Hex(), swapLog.To.Hex(),
				swapLog.Amount0In, swapLog.Amount1In, swapLog.Amount0Out, swapLog.Amount1Out, blockTime)

			if err != nil {
				return model.NewServerError().Err(fmt.Errorf("save swap event error: %w", err))
			}
		}

		return nil
	})
}

func NewSwapEventRepositoryRepository(pool *pgxpool.Pool) *SwapEventRepository {
	return &SwapEventRepository{
		BaseRepository: newBaseRepository(pool),
	}
}
