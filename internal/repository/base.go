package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"trading-ace/model"
)

type TransactionFunction = func(pgx.Tx) *model.AppError

type Repository interface {
	Pool() *pgxpool.Pool
}

type BaseRepository struct {
	p *pgxpool.Pool
}

func (r *BaseRepository) transaction(ctx context.Context, txFunc TransactionFunction) *model.AppError {
	tx, err := r.p.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return model.NewServerError().Err(fmt.Errorf("begin transaction error: %w", err))
	}

	appErr := txFunc(tx)

	if appErr != nil {
		_ = tx.Rollback(ctx)
		return appErr
	}

	err = tx.Commit(ctx)

	if err != nil && errors.Is(err, pgx.ErrTxCommitRollback) {
		return model.NewServerError().Err(fmt.Errorf("commit transaction error: %w", err))
	}

	return nil
}

func newBaseRepository(pool *pgxpool.Pool) *BaseRepository {
	return &BaseRepository{p: pool}
}
