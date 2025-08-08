package interfaces

import (
	"context"
	"testYTask/internal/domain/models"
)

type NexusI interface {
	Start() error
	Stop() error
}

type RegistryClientI interface {
	LoadValidCoins(ctx context.Context) (map[string]bool, error)
	CurrentData(ctx context.Context, coin string) (*models.CoinPrice, error)
}
