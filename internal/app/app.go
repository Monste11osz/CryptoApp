package app

import (
	"context"
	"testYTask/internal/config"
	"testYTask/internal/delivery/http"
	"testYTask/internal/delivery/http/server"
	"testYTask/internal/delivery/http/v1/handlers"
	"testYTask/internal/domain/interfaces"
	cli "testYTask/internal/http"
	"testYTask/internal/infrastructure/db"
	"testYTask/internal/infrastructure/db/repository"
	"testYTask/internal/usecase/job"

	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// App представляет основное приложение с его зависимостями.
type App struct {
	cfg       *config.Config
	nexus     interfaces.NexusI
	db        *pgxpool.Pool
	scheduler gocron.Scheduler
}

// NewApp создаёт новый экземпляр приложения с заданной конфигурацией.
//
// Параметры:
//   - cfg: указатель на объект конфигурации
//
// Возвращает:
//   - указатель на новый объект App.
func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

// Init инициализирует основные компоненты приложения.
//
// Параметры:
//   - ctx: контекст выполнения
//
// Возвращает ошибку, если инициализация не удалась.
func (a *App) Init(ctx context.Context) error {
	var err error

	// Инициализация логгера
	if err = InitLogger(a.cfg.App.Stage); err != nil {
		zap.L().Error("InitLogger failed", zap.Error(err))
		panic(err)
	}

	// Установка подключения к PostgreSQL
	a.db, err = db.Connection(ctx, a.cfg.Db)
	if err != nil {
		zap.L().Error("NewConnect failed", zap.Error(err))
		panic(err)
	}

	// Создание HTTP клиента для сервиса
	registryClient := cli.NewRegistryClient(a.cfg.Exchange)

	mapCoins, err := registryClient.LoadValidCoins(ctx)
	if err != nil {
		zap.L().Error("LoadValidCoins failed", zap.Error(err))
		panic(err)
	}

	// Инициализация репозиториев PostgreSQL
	majorRepository := repository.NewMajorRepository(a.db)
	jobRepository := repository.NewJobRepository(a.db)

	// Инициализация HTTP обработчиков
	majorHandler := handlers.NewMajorHandler(majorRepository, mapCoins)
	commonHandler := http.NewCommonHandler(a.db)

	// Инициализация задачи для загрузки
	uploadJob := job.NewUploadJob(
		jobRepository,
		registryClient,
	)

	//// Инициализация планировщика задач
	a.scheduler, err = InitScheduler(uploadJob)
	if err != nil {
		return err
	}

	// Создание и регистрация маршрутов HTTP-сервера
	navigator := server.NewNavigator(a.cfg)
	navigator.RegisterRoutes(commonHandler, majorHandler)

	a.nexus = server.NewNexus(a.cfg, navigator.Engine)
	return nil
}

// Run запускает планировщик, HTTP-сервер и Kafka-консьюмер.
//
// Функция запускает несколько горутин для параллельного выполнения и ожидает сигнал завершения.
func (a *App) Run(ctx context.Context) {
	defer a.onShutdown()

	// Запуск планировщика задач
	a.scheduler.Start()

	// Создание errgroup для параллельного выполнения горутин
	g, gCtx := errgroup.WithContext(ctx)

	// Запуск HTTP-сервера
	g.Go(func() error {
		return a.nexus.Start()
	})

	// Ожидание завершения работы по сигналу отмены контекста
	<-gCtx.Done()

	if err := g.Wait(); err != nil {
		zap.L().Error("Service shutdown error", zap.Error(err))
	}

	_ = zap.L().Sync()
}

// onShutdown выполняет завершающие действия: закрытие соединений и остановка планировщика.
//
// Функция вызывается при завершении работы приложения.
func (a *App) onShutdown() {
	a.closeConnection()
}

// closeConnections закрывает все открытые подключения.
func (a *App) closeConnection() {
	connectionPool := []interface{}{
		a.db,
	}

	for _, pool := range connectionPool {
		if pool == nil {
			continue
		}
		switch v := pool.(type) {
		case interfaces.CloserI:
			if err := v.Close(); err != nil {
				zap.L().Error("Failed to close connection", zap.Error(err))
			}
		}
	}
}
