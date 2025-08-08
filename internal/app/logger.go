package app

import (
	"strings"

	"go.uber.org/zap"
)

// InitLogger инициализирует глобальный логгер на основе стадии запуска приложения.
//
// Параметры:
//   - stage: стадия приложения ("production" или иное значение)
//
// Возвращает:
//   - ошибку, если не удалось инициализировать логгер.
func InitLogger(stage string) error {
	var (
		logger *zap.Logger
		err    error
	)

	if strings.EqualFold(stage, "production") {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	zap.L().Info("Successfully initialized logger")
	return nil
}
