package repository

import (
	"context"
	"errors"
	"math"
	"strings"
	"testYTask/internal/common"
	"testYTask/internal/domain/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type MajorRepository struct {
	db *pgxpool.Pool
}

func NewMajorRepository(db *pgxpool.Pool) *MajorRepository {
	return &MajorRepository{
		db: db,
	}
}

func (m *MajorRepository) AddingNewCoin(ctx context.Context, coin *models.Coin) error {
	dbCtx, cancel := context.WithTimeout(ctx, common.PostgresDBQueryTimeout)
	defer cancel()

	if _, err := m.db.Exec(dbCtx, queryAddingNewCoin, strings.ToLower(coin.NameCoin)); err != nil {
		zap.L().Error("Error adding new coin", zap.Error(err))
		return err
	}
	return nil
}

func (m *MajorRepository) DeleteCoin(ctx context.Context, coin *models.Coin) error {
	dbCtx, cancel := context.WithTimeout(ctx, common.PostgresDBQueryTimeout)
	defer cancel()

	result, err := m.db.Exec(dbCtx, queryDeleteCoin, coin.NameCoin)
	if err != nil {
		zap.L().Error("Error deleting coin", zap.Error(err))
	}

	rows := result.RowsAffected()
	if rows == common.Zero {
		return common.ErrCoinNotFound
	}

	return nil
}

func (m *MajorRepository) GetPrice(ctx context.Context, req *models.PriceRequest) (*models.PriceResponse, error) {
	dbCtx, cancel := context.WithTimeout(ctx, common.PostgresDBQueryTimeout)
	defer cancel()

	var DbResponse models.DbResponse

	row := m.db.QueryRow(dbCtx, queryGetPriceForCoin, strings.ToLower(req.Coin), req.Timestamp)
	err := row.Scan(&DbResponse.Coin, &DbResponse.Price, &DbResponse.Precision, &DbResponse.Currency, &DbResponse.Timestamp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrPriceNotFound
		}

		zap.L().Error("Error getting price", zap.Error(err))
		return nil, err
	}

	price := DbResponse.Price / math.Pow10(DbResponse.Precision)

	return &models.PriceResponse{
		Coin:      DbResponse.Coin,
		Price:     price,
		Currency:  DbResponse.Currency,
		Timestamp: DbResponse.Timestamp,
	}, nil
}
