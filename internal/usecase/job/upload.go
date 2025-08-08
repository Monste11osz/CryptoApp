package job

import (
	"context"
	"testYTask/internal/common"
	"testYTask/internal/domain/interfaces"

	"go.uber.org/zap"
)

// UploadJob представляет задачу запроса на сторонний источник и дальнейшим сохранением данных.
type UploadJob struct {
	jobRepository  interfaces.JobRepositoryI
	registryClient interfaces.RegistryClientI
}

// NewUploadJob создает новую задачу для запроса на сторонний источник и дальнейшим сохранением данных.
//
// Параметры:
//   - registryClient: клиент для работы с внешними HTTP-сервисами
//   - jobRepository: репозиторий для работы с Postgres
//
// Возвращает указатель на созданную задачу.
func NewUploadJob(jobRepository interfaces.JobRepositoryI, registryClient interfaces.RegistryClientI) *UploadJob {
	return &UploadJob{
		jobRepository:  jobRepository,
		registryClient: registryClient,
	}
}

// Run запускает задачу запроса и сохранения данных.
func (job *UploadJob) Run() {
	zap.L().Info("Starting upload job")

	ctx := context.Background()

	listCoins, err := job.jobRepository.ListOfCurrentCoins(ctx)
	if err != nil {
		zap.L().Error("ListOfCurrentCoins failed", zap.Error(err))
		return
	}

	for _, coin := range listCoins {

		data, err := job.registryClient.CurrentData(ctx, coin)
		if err != nil {
			zap.L().Error("data retrieval error", zap.Error(err), zap.String("name:", coin))
			continue
		}
		if err := job.jobRepository.CoinDataUpdate(ctx, coin, data, common.Precision, common.DefaultCurrency); err != nil {
			zap.L().Error("update coin data error", zap.Error(err), zap.String("name:", coin))
			continue
		}

	}

	zap.L().Info("Finished upload job")
}
