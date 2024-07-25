package worker

import (
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"time"
	"trading-ace/model"
)

func (a *App) startLogUpdater(logStates []model.AddressLogState) {
	for _, logState := range logStates {
		go func(logState model.AddressLogState) {
			a.updateLogs(&logState)
		}(logState)
	}
}

func (a *App) updateLogs(logState *model.AddressLogState) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress(logState.ContractAddress),
		},
	}

	currentBlock := logState.CurrentBlockNumber
	topics := make([][]common.Hash, 4)
	hasTopic := false

	for i, topic := range []sql.NullString{logState.Topic0, logState.Topic1, logState.Topic2, logState.Topic3} {
		if topic.Valid {
			hasTopic = true
			topics[i] = []common.Hash{common.HexToHash(topic.String)}
		}
	}

	if hasTopic {
		query.Topics = topics
	}

	ticker := time.NewTicker(time.Second * time.Duration(a.cfg.AddressLog.UpdateInterval))

	for {
		select {
		case <-ticker.C:
			blockNumber, err := a.rpcClient.BlockNumber(a.ctx)
			if err != nil {
				log.Println("get block number error: ", err)
				continue
			}
			maxBlock := blockNumber - a.cfg.AddressLog.ConfirmationBlocks
			fromBlock := currentBlock

		loopBlocks:
			for fromBlock < maxBlock {
				toBlock := fromBlock + a.cfg.AddressLog.BatchBlocks
				log.Printf("update logs contract address: %s, from block: %d, to block; %d\n", logState.ContractAddress, fromBlock, toBlock)
				if toBlock > maxBlock {
					toBlock = maxBlock
				}

				query.FromBlock = new(big.Int).SetUint64(fromBlock)
				query.ToBlock = new(big.Int).SetUint64(toBlock)

				appErr := a.updateLogsByQuery(query)
				if appErr != nil {
					log.Printf("contract address: %s, error: %v\n", logState.ContractAddress, appErr)
					break loopBlocks
				}

				var blockTime time.Time
				blockTime, appErr = a.getOrUpdateBlockTime(blockNumber)
				if appErr != nil {
					log.Println(appErr)
					break loopBlocks
				}

				logState.CurrentBlockNumber = toBlock
				logState.CurrentBlockTime = blockTime
				logState.UpdatedAt = time.Now()

				appErr = a.addressLogStateRepository.UpdateCurrentBlock(a.ctx, *logState)
				if appErr != nil {
					log.Printf("contract address: %s, error: %v\n", logState.ContractAddress, appErr)
					break loopBlocks
				}

				fromBlock = toBlock
				a.addressLogStateChan <- *logState
			}

			if fromBlock > currentBlock {
				currentBlock = fromBlock
			}
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *App) updateLogsByQuery(query ethereum.FilterQuery) *model.AppError {
	a.group.Add(1)

	defer func() {
		a.group.Done()
	}()

	logs, err := a.rpcClient.FilterLogs(a.ctx, query)
	if err != nil {
		return model.NewAppError(model.OperateFailed).Err(fmt.Errorf("eth get logs error: %w", err))
	}

	var blockNumbers []uint64
	blockNumberAndLogsMap := map[uint64][]types.Log{}

	for _, item := range logs {
		blockNumbers = append(blockNumbers, item.BlockNumber)
		blockNumberAndLogsMap[item.BlockNumber] = append(blockNumberAndLogsMap[item.BlockNumber], item)
	}

	for blockNumber, blockLogs := range blockNumberAndLogsMap {
		blockTime, appErr := a.getOrUpdateBlockTime(blockNumber)
		if appErr != nil {
			return appErr
		}

		appErr = a.saveSwapEvents(blockTime, blockLogs)
		if appErr != nil {
			log.Println(appErr)
			return appErr
		}
	}

	return nil
}
