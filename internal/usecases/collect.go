package usecases

import (
	"context"
	"metric-server/internal/adapters/http/api_v01/dto"
	"metric-server/internal/models"
)

type Storage interface {
	Save(ctx context.Context, metric models.Metric) error
	SaveMany(ctx context.Context, metrics []models.Metric) error
	FindByNameAndType(ctx context.Context, name string, metricType string) (models.Metric, error)
	FindAllByType(ctx context.Context, metricType string) ([]models.Metric, error)
}

type Mapper interface {
	ToMetricModel(ctx context.Context, dto dto.CollectDto) (models.Metric, error)
	ToMetricModelList(ctx context.Context, dto []dto.CollectDto) ([]models.Metric, error)
}

type CollectUseCase struct {
	s Storage
	m Mapper
}

func NewCollectUseCase(s Storage, m Mapper) *CollectUseCase {
	return &CollectUseCase{
		s: s,
		m: m,
	}
}

func (c CollectUseCase) CollectOne(ctx context.Context, dto dto.CollectDto) error {
	model, err := c.m.ToMetricModel(ctx, dto)
	if err != nil {
		return err
	}

	return c.s.Save(ctx, model)
}

func (c CollectUseCase) CollectMany(ctx context.Context, metrics []dto.CollectDto) error {
	metricModels, err := c.m.ToMetricModelList(ctx, metrics)
	if err != nil {
		return err
	}

	return c.s.SaveMany(ctx, metricModels)
}
