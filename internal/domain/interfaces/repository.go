package interfaces

import (
	"context"
	"testYTask/internal/domain/models"
)

type MajorRepositoryI interface {
	AddingNewCoin(ctx context.Context, coin *models.Coin) error
	DeleteCoin(ctx context.Context, coin *models.Coin) error
	GetPrice(ctx context.Context, coin *models.PriceRequest) (*models.PriceResponse, error)
}

type JobRepositoryI interface {
	ListOfCurrentCoins(ctx context.Context) ([]string, error)
	CoinDataUpdate(ctx context.Context, coin string, data *models.CoinPrice, precision int, currency string) error
}
