package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"trading-ace/model"
)

type BlockRepository struct {
	*BaseRepository
}

func (r *BlockRepository) GetBlockTime(ctx context.Context, blockNumber uint64) (time.Time, *model.AppError) {
	row := r.p.QueryRow(ctx, `
SELECT block_timestamp
FROM block_info
WHERE id = $1
`, blockNumber)

	var blockTime time.Time

	err := row.Scan(&blockTime)
	if err != nil {
		return blockTime, model.NewAppError(model.OperateFailed).Err(fmt.Errorf("get block info error: %w", err))
	}

	return blockTime, nil
}

func (r *BlockRepository) Create(ctx context.Context, blockNumber uint64, blockTime time.Time) *model.AppError {
	_, err := r.p.Exec(ctx, `
INSERT INTO block_info (id, block_timestamp, updated_at)
VALUES ($1, $2, NOW())
ON CONFLICT (id) DO NOTHING
`, blockNumber, blockTime)
	if err != nil {
		return model.NewAppError(model.OperateFailed).Err(fmt.Errorf("save block info error: %w", err))
	}

	return nil
}

func NewBlockRepository(pool *pgxpool.Pool) *BlockRepository {
	return &BlockRepository{
		BaseRepository: newBaseRepository(pool),
	}
}
