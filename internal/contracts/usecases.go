package contracts

import (
	"context"
	"github.com/0db0/metric-server/internal/dto"
	"github.com/0db0/metric-server/internal/models"
)

//go:generate mockgen -source=usecases.go -package=contracts -destination=../mocks/contracts/mock_usecases.go
type CollectUseCase interface {
	CollectOne(ctx context.Context, metric dto.CollectDto) error
	CollectMany(ctx context.Context, metrics []dto.CollectDto) error
}

type GiveUseCase interface {
	GetValue(ctx context.Context, dto dto.ValueDto) (models.Metric, error)
}
