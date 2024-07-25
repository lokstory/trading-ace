package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type BlockNumberResponse struct {
	Height    uint64 `json:"height"`
	Timestamp uint64 `json:"timestamp"`
}

type UserPointTaskStateWithTaskInfo struct {
	Address         string           `json:"address"`
	Point           uint64           `json:"point"`
	Volume          *decimal.Decimal `json:"volume"`
	Completed       bool             `json:"completed"`
	UpdatedAt       *time.Time       `json:"updated_at"`
	TaskName        string           `json:"task_name"`
	TaskType        string           `json:"task_type"`
	ContractAddress string           `json:"contract_address"`
	TokenAddress    string           `json:"token_address"`
	TaskVolume      *decimal.Decimal `json:"task_volume"`
	TaskRewardPoint uint64           `json:"task_reward_point"`
	TaskStatus      string           `json:"task_status"`
	StartTime       time.Time        `json:"start_time"`
	EndTime         time.Time        `json:"end_time"`
}

type UserPointHistoryWithTaskInfo struct {
	Address   string    `json:"address"`
	Point     uint64    `json:"point"`
	TaskName  string    `json:"task_name"`
	TaskType  string    `json:"task_type"`
	CreatedAt time.Time `json:"created_at"`
}
