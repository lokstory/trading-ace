package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"trading-ace/model"
)

type UserPointTaskStateRepository struct {
	*BaseRepository
}

func (r *UserPointTaskStateRepository) ListWithTaskByAddress(ctx context.Context, address string) ([]model.UserPointTaskStateWithTaskInfo, *model.AppError) {
	rows, err := r.p.Query(ctx, `
SELECT COALESCE(s.address, $1)                                           AS address,
       COALESCE(s.point, 0)                                              AS point,
       COALESCE(s.volume, 0)                                             AS volume,
       CASE
         WHEN s.status = 'COMPLETED' THEN TRUE
         ELSE FALSE
         END                                                             AS completed,
       s.updated_at,
       t.name                                                            AS task_name,
       t.task_type                                                       AS task_type,
       t.contract_address                                                AS contract_address,
       t.token_address                                                   AS token_address,
       t.volume                                                          AS task_volume,
       t.reward_point                                                    AS task_reward_point,
       t.start_time                                                      AS start_time,
       t.end_time                                                        AS end_time,
       t.status                                                          AS task_status
FROM point_task AS t
     LEFT JOIN user_point_task_state AS s
               ON s.point_task_id = t.id AND s.address = $1
ORDER BY t.id
`, address)

	if err != nil {
		return nil, model.NewServerError().Err(fmt.Errorf("query user point task state error: %w", err))
	}

	defer rows.Close()

	var list []model.UserPointTaskStateWithTaskInfo

	for rows.Next() {
		var item model.UserPointTaskStateWithTaskInfo
		if err := rows.Scan(&item.Address, &item.Point, &item.Volume, &item.Completed, &item.UpdatedAt, &item.TaskName,
			&item.TaskType, &item.ContractAddress, &item.TokenAddress, &item.TaskVolume, &item.TaskRewardPoint, &item.StartTime, &item.EndTime, &item.TaskStatus); err != nil {
			return nil, model.NewServerError().Err(fmt.Errorf("scan user point task state error: %w", err))
		}
		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, model.NewServerError().Err(fmt.Errorf("read user point task state error: %w", err))
	}

	return list, nil
}

func NewUserPointTaskStateRepository(pool *pgxpool.Pool) *UserPointTaskStateRepository {
	return &UserPointTaskStateRepository{
		BaseRepository: newBaseRepository(pool),
	}
}
