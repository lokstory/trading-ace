package worker

import (
	"github.com/shopspring/decimal"
	"log"
	"trading-ace/model"
)

func (a *App) initPointTasks() *model.AppError {
	onBoardingVolume := a.cfg.Task.OnBoarding.Volume

	accountTradingVolumeTask := &model.PointTask{
		Name:            "Onboarding Task",
		TaskType:        model.TaskAccountTradingVolume,
		ContractAddress: a.cfg.Task.ContractAddress,
		TokenAddress:    a.cfg.Task.TokenAddress,
		Volume:          &onBoardingVolume,
		RewardPoint:     a.cfg.Task.OnBoarding.Point,
		SettlementType:  model.SettlementInstant,
		Status:          model.StatusCreated,
		StartTime:       a.cfg.Task.StartTime,
		EndTime:         a.cfg.Task.EndTime,
	}

	zero := decimal.Zero
	sharedPoolTask := &model.PointTask{
		Name:            "Share Pool Task",
		TaskType:        model.TaskSharePool,
		ContractAddress: a.cfg.Task.ContractAddress,
		TokenAddress:    a.cfg.Task.TokenAddress,
		Volume:          &zero,
		RewardPoint:     a.cfg.Task.SharePool.Point,
		SettlementType:  model.SettlementWeekly,
		Status:          model.StatusCreated,
		StartTime:       a.cfg.Task.StartTime,
		EndTime:         a.cfg.Task.EndTime,
	}

	appErr := a.pointTaskRepository.CreateTasks(a.ctx, []*model.PointTask{accountTradingVolumeTask, sharedPoolTask})
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a *App) startPointTaskUpdater(tasks []model.PointTask) {
	var refreshTasks = func() {
		list, appErr := a.pointTaskRepository.ListByStatus(a.ctx, model.StatusCreated)
		if appErr != nil {
			log.Println(appErr)
		} else {
			tasks = list
		}
	}

	for {
		select {
		case item := <-a.addressLogStateChan:
			blockTime := item.CurrentBlockTime

			for _, item := range tasks {
				if blockTime.Before(item.StartTime) {
					continue
				}

				switch item.TaskType {
				case model.TaskAccountTradingVolume:
					apiErr := a.pointTaskRepository.UpdateAccountTradingVolume(a.ctx, item.ID)
					if apiErr != nil {
						log.Println(apiErr)
					}
				case model.TaskSharePool:
					if blockTime.Before(item.EndTime) {
						continue
					}

					apiErr := a.pointTaskRepository.UpdateSharePool(a.ctx, item.ID)
					if apiErr != nil {
						log.Println(apiErr)
					}
				}
			}

			refreshTasks()
		case <-a.ctx.Done():
			return
		}
	}
}
