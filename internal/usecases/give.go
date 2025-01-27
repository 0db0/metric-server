package usecases

import (
	"context"
	"github.com/0db0/metric-server/internal/dto"
	"github.com/0db0/metric-server/internal/models"
)

type GiveUseCase struct {
	s Storage
}

func NewGiveUseCase(s Storage) *GiveUseCase {
	return &GiveUseCase{
		s: s,
	}
}

func (g GiveUseCase) GetValue(ctx context.Context, dto dto.ValueDto) (models.Metric, error) {
	metric, err := g.s.FindByNameAndType(ctx, dto.ID, dto.MType)

	if err != nil {
		return models.Metric{}, err
	}

	return metric, nil
}
