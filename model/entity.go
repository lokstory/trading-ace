package model

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"time"
)

const (
	TaskAccountTradingVolume = "ACCOUNT_TRADING_VOLUME"
	TaskSharePool            = "SHARE_POOL"
)

const (
	SettlementInstant = "INSTANT"
	SettlementWeekly  = "WEEKLY"
)

const (
	StatusCreated   = "CREATED"
	StatusFinished  = "FINISHED"
	StatusCompleted = "COMPLETED"
)

type AddressLogState struct {
	ID                 uint64         `json:"id"`
	ContractAddress    string         `json:"contract_address"`
	Topic0             sql.NullString `json:"topic0"`
	Topic1             sql.NullString `json:"topic1"`
	Topic2             sql.NullString `json:"topic2"`
	Topic3             sql.NullString `json:"topic3"`
	StartBlockNumber   uint64         `json:"start_block_number"`
	CurrentBlockNumber uint64         `json:"current_block_number"`
	CurrentBlockTime   time.Time      `json:"current_block_time"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

type PointTask struct {
	ID                 uint64           `json:"id"`
	Name               string           `json:"name"`
	TaskType           string           `json:"task_type"`
	ContractAddress    string           `json:"contract_address"`
	TokenAddress       string           `json:"token_address"`
	Volume             *decimal.Decimal `json:"volume"`
	RewardPoint        uint64           `json:"reward_point"`
	SettlementType     string           `json:"settlement_type"`
	Status             string           `json:"status"`
	StartTime          time.Time        `json:"start_time"`
	EndTime            time.Time        `json:"end_time"`
	UpdatedAtBlockTime time.Time        `json:"updated_at_block_time"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}
