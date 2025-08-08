package repository

import (
	"context"
	"math"
	"testYTask/internal/common"
	"testYTask/internal/domain/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
	return &JobRepository{
		db: db,
	}
}

func (r *JobRepository) CoinDataUpdate(ctx context.Context, coin string, data *models.CoinPrice, precision int, currency string) error {
	dbCtx, cancel := context.WithTimeout(ctx, common.PostgresDBQueryTimeout)
	defer cancel()

	timestamp := time.Now().Unix()

	multiplier := math.Pow10(precision)
	intPrice := int64(data.USD * multiplier)

	if _, err := r.db.Exec(dbCtx, queryUpdatePriceCoin, coin, timestamp, intPrice, precision, currency); err != nil {
		zap.L().Error("update coin data error", zap.Error(err), zap.String("name:", coin))
		return err
	}
	return nil
}

func (r *JobRepository) ListOfCurrentCoins(ctx context.Context) ([]string, error) {
	dbCtx, cancel := context.WithTimeout(ctx, common.PostgresDBQueryTimeout)
	defer cancel()

	rows, err := r.db.Query(dbCtx, queryListRequest)
	if err != nil {
		zap.L().Error("failed to list current coins", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var coins []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			zap.L().Error("Scan error", zap.Error(err))
			return nil, err
		}
		coins = append(coins, c)
	}
	return coins, nil

}
