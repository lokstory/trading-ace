package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"trading-ace/model"
)

type UserPointHistoryRepository struct {
	*BaseRepository
}

func (r *UserPointHistoryRepository) ListWithTaskByAddress(ctx context.Context, address string) ([]model.UserPointHistoryWithTaskInfo, *model.AppError) {
	rows, err := r.p.Query(ctx, `
SELECT h.address,
       h.point,
       h.created_at,
       t.name             AS task_name,
       t.task_type        AS task_type
FROM user_point_history AS h
     INNER JOIN point_task AS t ON t.id = h.point_task_id
WHERE h.address = $1
`, address)

	if err != nil {
		return nil, model.NewServerError().Err(fmt.Errorf("query user point history error: %w", err))
	}

	defer rows.Close()

	var list []model.UserPointHistoryWithTaskInfo

	for rows.Next() {
		var item model.UserPointHistoryWithTaskInfo
		if err := rows.Scan(&item.Address, &item.Point, &item.CreatedAt,
			&item.TaskName, &item.TaskType); err != nil {
			return nil, model.NewServerError().Err(fmt.Errorf("scan user point history error: %w", err))
		}
		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, model.NewServerError().Err(fmt.Errorf("read user point history error: %w", err))
	}

	return list, nil
}

func NewUserPointHistoryRepository(pool *pgxpool.Pool) *UserPointHistoryRepository {
	return &UserPointHistoryRepository{
		BaseRepository: newBaseRepository(pool),
	}
}
