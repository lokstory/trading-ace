package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"trading-ace/model"
)

type PointTaskRepository struct {
	*BaseRepository
}

func (r *PointTaskRepository) ListAll(ctx context.Context) ([]model.PointTask, *model.AppError) {
	rows, err := r.p.Query(ctx, `
SELECT id,
       name,
       task_type,
       contract_address,
       volume,
       reward_point,
       settlement_type,
       status,
       start_time,
       end_time,
       updated_at_block_time,
       created_at,
       updated_at
FROM public.point_task
`)

	if err != nil {
		return nil, model.NewServerError().Err(fmt.Errorf("query point task error: %w", err))
	}

	defer rows.Close()

	return r.decodeRows(rows)
}

func (r *PointTaskRepository) ListByStatus(ctx context.Context, taskStatus string) ([]model.PointTask, *model.AppError) {
	rows, err := r.p.Query(ctx, `
SELECT id,
       name,
       task_type,
       contract_address,
       volume,
       reward_point,
       settlement_type,
       status,
       start_time,
       end_time,
       updated_at_block_time,
       created_at,
       updated_at
FROM public.point_task
WHERE status = $1
`, taskStatus)

	if err != nil {
		return nil, model.NewServerError().Err(fmt.Errorf("query point task by status error: %w", err))
	}

	defer rows.Close()

	return r.decodeRows(rows)
}

func (r *PointTaskRepository) CreateTasks(ctx context.Context, items []*model.PointTask) *model.AppError {
	return r.transaction(ctx, func(tx pgx.Tx) *model.AppError {
		for _, item := range items {
			_, err := tx.Exec(ctx, `
INSERT INTO public.point_task (name, task_type, contract_address, token_address, volume, reward_point, settlement_type, status, start_time, end_time, updated_at_block_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`, item.Name, item.TaskType, item.ContractAddress, item.TokenAddress, item.Volume, item.RewardPoint, item.SettlementType, item.Status, item.StartTime, item.EndTime, time.Unix(0, 0))

			if err != nil {
				return model.NewServerError().Err(fmt.Errorf("create point tasks error: %w", err))
			}
		}

		return nil
	})
}

func (r *PointTaskRepository) UpdateAccountTradingVolume(ctx context.Context, id uint64) *model.AppError {
	return r.transaction(ctx, func(tx pgx.Tx) *model.AppError {
		_, err := tx.Exec(ctx, `
CALL public.update_account_trading_volume_point_task($1)
`, id)

		if err != nil {
			return model.NewServerError().Err(fmt.Errorf("update account trading volume point task error: %w", err))
		}

		return nil
	})
}

func (r *PointTaskRepository) UpdateSharePool(ctx context.Context, id uint64) *model.AppError {
	return r.transaction(ctx, func(tx pgx.Tx) *model.AppError {
		_, err := tx.Exec(ctx, `
CALL public.update_share_pool_point_task($1)
`, id)

		if err != nil {
			return model.NewServerError().Err(fmt.Errorf("update share pool point task error: %w", err))
		}

		return nil
	})
}

func (r *PointTaskRepository) decodeRows(rows pgx.Rows) ([]model.PointTask, *model.AppError) {
	var list []model.PointTask

	for rows.Next() {
		var item model.PointTask
		if err := rows.Scan(&item.ID, &item.Name, &item.TaskType, &item.ContractAddress, &item.Volume, &item.RewardPoint,
			&item.SettlementType, &item.Status, &item.StartTime, &item.EndTime, &item.UpdatedAtBlockTime, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, model.NewServerError().Err(fmt.Errorf("scan point task error: %w", err))
		}
		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, model.NewServerError().Err(fmt.Errorf("read point task error: %w", err))
	}

	return list, nil
}

func NewTaskRepository(pool *pgxpool.Pool) *PointTaskRepository {
	return &PointTaskRepository{
		BaseRepository: newBaseRepository(pool),
	}
}
