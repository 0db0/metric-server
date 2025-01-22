package api_v01

import (
	"context"
	"metric-server/internal/models"
	"metric-server/internal/pkg/logger"
)

type CollectUseCase interface {
	CollectOne(ctx context.Context, metric CollectDto) error
	CollectMany(ctx context.Context, metrics []CollectDto) error
}

type GiveUseCase interface {
	GetValue(ctx context.Context, dto ValueDto) (models.Metric, error)
}

type MetricAdapter struct {
	c   CollectUseCase
	g   GiveUseCase
	log logger.Interface
}

func NewMetricAdapter(c CollectUseCase, g GiveUseCase, log logger.Interface) *MetricAdapter {
	return &MetricAdapter{
		c:   c,
		g:   g,
		log: log,
	}
}
