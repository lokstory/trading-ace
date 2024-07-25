package api

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"trading-ace/internal/database"
	"trading-ace/internal/repository"
	"trading-ace/internal/util"
	"trading-ace/model"
)

type Cfg struct {
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	Http struct {
		Address string `yaml:"address"`
	} `yaml:"http"`
}

type App struct {
	ctx                          context.Context
	cfg                          *Cfg
	f                            *fiber.App
	dbPool                       *pgxpool.Pool
	userPointTaskStateRepository *repository.UserPointTaskStateRepository
	userPointHistoryRepository   *repository.UserPointHistoryRepository
}

func (a *App) Start() {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			var appErr *model.AppError

			log.Println("request error:", err)

			if errors.As(err, &appErr) {

				ctx.Status(appErr.HTTPStatusCode)

				_ = ctx.JSON(appErr)
				return nil
			}

			_ = ctx.Status(http.StatusInternalServerError).JSON(model.NewServerError())

			return nil
		},
	})

	a.f = app

	app.Use(panicRecover())
	a.initRoute(app)

	app.Use(func(c fiber.Ctx) error {
		_ = c.Status(http.StatusNotFound).JSON(model.NewAppError(model.NotFound))
		return nil
	})

	err := app.Listen(a.cfg.Http.Address)
	if err != nil {
		log.Println(err)
	}
}

func (a *App) Close() {
	err := a.f.Shutdown()
	if err != nil {
		log.Println("shutdown http error:", err)
	}
	a.dbPool.Close()
}

func NewApp(ctx context.Context, cfgPath string) (*App, error) {
	cfg := &Cfg{}
	err := util.LoadYAML(cfgPath, cfg)
	if err != nil {
		return nil, err
	}

	dbPool, err := database.NewDB(cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:                          ctx,
		cfg:                          cfg,
		dbPool:                       dbPool,
		userPointTaskStateRepository: repository.NewUserPointTaskStateRepository(dbPool),
		userPointHistoryRepository:   repository.NewUserPointHistoryRepository(dbPool),
	}, nil
}
