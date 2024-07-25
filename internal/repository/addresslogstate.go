package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"trading-ace/model"
)

type AddressLogStateRepository struct {
	*BaseRepository
}

func (r *AddressLogStateRepository) ListAll(ctx context.Context) ([]model.AddressLogState, *model.AppError) {
	rows, err := r.p.Query(ctx, `
SELECT id, contract_address, topic0, topic1, topic2, topic3, start_block_number, current_block_number, created_at, updated_at
FROM address_log_state
`)

	if err != nil {
		return nil, model.NewServerError().Err(fmt.Errorf("query address log state error: %w", err))
	}

	defer rows.Close()

	var list []model.AddressLogState

	for rows.Next() {
		var item model.AddressLogState
		if err := rows.Scan(&item.ID, &item.ContractAddress, &item.Topic0, &item.Topic1, &item.Topic2, &item.Topic3,
			&item.StartBlockNumber, &item.CurrentBlockNumber, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, model.NewServerError().Err(fmt.Errorf("scan address log state error: %w", err))
		}
		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		return nil, model.NewServerError().Err(fmt.Errorf("read address log state error: %w", err))
	}

	return list, nil
}

func (r *AddressLogStateRepository) CreateOrUpdate(ctx context.Context, item *model.AddressLogState) *model.AppError {
	_, err := r.p.Exec(ctx, `
INSERT INTO public.address_log_state (contract_address, topic0, topic1, topic2, topic3, start_block_number, current_block_number,
                                      current_block_time, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (contract_address) DO UPDATE
  SET topic0               = $2,
      topic1               = $3,
      topic2               = $4,
      topic3               = $5,
      start_block_number   = $6,
      current_block_number = $7,
      current_block_time = $8,
      created_at = $9,
      updated_at = $10
`, item.ContractAddress, item.Topic0, item.Topic1, item.Topic2, item.Topic3, item.StartBlockNumber, item.CurrentBlockNumber, item.CurrentBlockTime, item.CreatedAt, item.UpdatedAt)

	if err != nil {
		return model.NewServerError().Err(fmt.Errorf("upsert address log state error: %w", err))
	}

	return nil
}

func (r *AddressLogStateRepository) UpdateCurrentBlock(ctx context.Context, logState model.AddressLogState) *model.AppError {
	ret, err := r.p.Exec(ctx, `
UPDATE public.address_log_state
SET current_block_number = $1, current_block_time = $2, updated_at = $3 WHERE id = $4
`, logState.CurrentBlockNumber, logState.CurrentBlockTime, logState.UpdatedAt, logState.ID)

	if err != nil {
		return model.NewServerError().Err(fmt.Errorf("update address log state error: %w", err))
	}

	if ret.RowsAffected() == 0 {
		return model.NewAppError(model.OperateFailed).Err(fmt.Errorf("update address log state failed, row affected: %d", ret.RowsAffected()))
	}

	return nil
}

func NewAddressLogStateRepository(pool *pgxpool.Pool) *AddressLogStateRepository {
	return &AddressLogStateRepository{
		BaseRepository: newBaseRepository(pool),
	}
}
