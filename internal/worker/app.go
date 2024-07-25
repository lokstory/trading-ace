package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"sync"
	"time"
	"trading-ace/internal/database"
	"trading-ace/internal/repository"
	"trading-ace/internal/util"
	"trading-ace/model"
)

type Cfg struct {
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	Blockchain struct {
		URL string `yaml:"url"`
	}
	AddressLog struct {
		UpdateInterval     uint64 `yaml:"update_interval"`
		ConfirmationBlocks uint64 `yaml:"confirmation_blocks"`
		BatchBlocks        uint64 `yaml:"batch_blocks"`
	} `yaml:"address_log"`
	Task struct {
		ContractAddress string    `yaml:"contract_address"`
		TokenAddress    string    `yaml:"token_address"`
		StartTime       time.Time `yaml:"start_time"`
		EndTime         time.Time `yaml:"end_time"`
		OnBoarding      struct {
			Volume decimal.Decimal `yaml:"volume"`
			Point  uint64          `yaml:"point"`
		} `yaml:"on_boarding"`
		SharePool struct {
			Point uint64 `yaml:"point"`
		} `yaml:"share_pool"`
	} `yaml:"task"`
}

type App struct {
	ctx                       context.Context
	cfg                       *Cfg
	group                     *sync.WaitGroup
	rpcClient                 *ethclient.Client
	dbPool                    *pgxpool.Pool
	addressLogStateChan       chan model.AddressLogState
	blockRepository           *repository.BlockRepository
	addressLogStateRepository *repository.AddressLogStateRepository
	swapEventRepository       *repository.SwapEventRepository
	pointTaskRepository       *repository.PointTaskRepository
}

func (a *App) Start() error {
	appErr := a.initData()
	if appErr != nil {
		return appErr
	}

	var logStates []model.AddressLogState
	var createdTasks []model.PointTask

	logStates, appErr = a.addressLogStateRepository.ListAll(a.ctx)
	if appErr != nil {
		return appErr
	}

	createdTasks, appErr = a.pointTaskRepository.ListByStatus(a.ctx, model.StatusCreated)
	if appErr != nil {
		return appErr
	}

	a.startLogUpdater(logStates)
	go a.startPointTaskUpdater(createdTasks)

	return nil
}

func (a *App) Close() {
	a.group.Wait()
	a.dbPool.Close()
}

func (a *App) initData() *model.AppError {
	tasks, appErr := a.pointTaskRepository.ListAll(a.ctx)
	if appErr != nil {
		return appErr
	}

	if len(tasks) > 0 {
		return nil
	}

	latestBlock, err := a.rpcClient.BlockByNumber(a.ctx, nil)
	if err != nil {
		return model.NewAppError(model.OperateFailed).Err(fmt.Errorf("get latest block error: %w", err))
	}

	var startBlockNumber uint64
	var startBlockTime time.Time

	if latestBlock.Time() <= uint64(a.cfg.Task.StartTime.Unix()) {
		startBlockNumber = latestBlock.NumberU64()
		startBlockTime = time.Unix(int64(latestBlock.Time()), 0)
	} else {
		blockNumberResponse, appErr := util.BlockNumberByTime(a.cfg.Task.StartTime)
		if appErr != nil {
			return appErr
		}

		startBlockNumber = blockNumberResponse.Height
		startBlockTime = time.Unix(int64(blockNumberResponse.Timestamp), 0)
	}

	now := time.Now()
	addressLogState := &model.AddressLogState{
		ContractAddress: a.cfg.Task.ContractAddress,
		Topic0: sql.NullString{
			Valid:  true,
			String: "0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822",
		},
		StartBlockNumber:   startBlockNumber,
		CurrentBlockNumber: startBlockNumber,
		CurrentBlockTime:   startBlockTime,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	appErr = a.addressLogStateRepository.CreateOrUpdate(a.ctx, addressLogState)
	if appErr != nil {
		return appErr
	}

	appErr = a.initPointTasks()
	return appErr
}

func NewApp(ctx context.Context, cfgPath string) (*App, error) {
	cfg := &Cfg{}
	err := util.LoadYAML(cfgPath, cfg)
	if err != nil {
		return nil, err
	}

	if cfg.AddressLog.BatchBlocks == 0 {
		return nil, errors.New("invalid batch blocks")
	}

	dbPool, err := database.NewDB(cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(cfg.Blockchain.URL)
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:                       ctx,
		cfg:                       cfg,
		group:                     &sync.WaitGroup{},
		rpcClient:                 client,
		dbPool:                    dbPool,
		addressLogStateChan:       make(chan model.AddressLogState, 32),
		blockRepository:           repository.NewBlockRepository(dbPool),
		addressLogStateRepository: repository.NewAddressLogStateRepository(dbPool),
		swapEventRepository:       repository.NewSwapEventRepositoryRepository(dbPool),
		pointTaskRepository:       repository.NewTaskRepository(dbPool),
	}, nil
}
