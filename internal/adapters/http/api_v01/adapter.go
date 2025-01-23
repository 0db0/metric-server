package api_v01

import (
	"context"
	"metric-server/internal/adapters/http/api_v01/dto"
	"metric-server/internal/models"
	"metric-server/internal/pkg/logger"
	"net/http"
)

//go:generate mockgen -source=adapter.go -package=api_v01 -destination=../../../mocks/adapters/http/api_v01/mock_usecases.go
type CollectUseCase interface {
	CollectOne(ctx context.Context, metric dto.CollectDto) error
	CollectMany(ctx context.Context, metrics []dto.CollectDto) error
}

type GiveUseCase interface {
	GetValue(ctx context.Context, dto dto.ValueDto) (models.Metric, error)
}

type RequestDtoBuilder interface {
	CreateCollectDto(r *http.Request) (dto.CollectDto, error)
}

type MetricAdapter struct {
	c   CollectUseCase
	g   GiveUseCase
	rb  RequestDtoBuilder
	log logger.Interface
}

func NewMetricAdapter(
	c CollectUseCase,
	g GiveUseCase,
	rb RequestDtoBuilder,
	log logger.Interface,
) *MetricAdapter {
	return &MetricAdapter{
		c:   c,
		g:   g,
		rb:  rb,
		log: log,
	}
}
