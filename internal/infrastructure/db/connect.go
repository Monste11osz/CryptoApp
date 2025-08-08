package db

import (
	"context"
	"fmt"
	"net/url"
	"testYTask/internal/common"
	"testYTask/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Connection устанавливает соединение с базой данных PostgreSQL,
// используя предоставленную конфигурацию. Функция создает пул соединений,
// проверяет его доступность и возвращает пул для дальнейшей работы.
//
// Параметры:
//   - ctx: контекст, используемый для управления временем ожидания подключения.
//   - cfg: конфигурация подключения к PostgreSQL.
//
// Возвращает:
//   - *pgxpool.Pool: пул соединений, если соединение успешно установлено.
//   - error: ошибку, если возникли проблемы при установлении соединения или проверке доступности.
func Connection(ctx context.Context, cfg *config.DbConf) (*pgxpool.Pool, error) {
	u := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(cfg.User, cfg.Pass),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.Name,
	}
	q := u.Query()
	u.RawQuery = q.Encode()
	dsn := u.String()
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		zap.L().Error("Error parsing pool configuration", zap.Error(err))
		return nil, err
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, common.ConnTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctxWithTimeout, poolConfig)
	if err != nil {
		zap.L().Error("Error creating connection pool", zap.Error(err))
		return nil, err
	}

	if err = pool.Ping(ctxWithTimeout); err != nil {
		zap.L().Error("Error pinging PostgreSQL database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Successfully connected to PostgreSQL")
	return pool, nil
}
