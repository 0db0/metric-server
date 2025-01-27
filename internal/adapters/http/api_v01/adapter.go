package api_v01

import (
	"metric-server/internal/contracts"
	"metric-server/internal/dto"
	"metric-server/internal/pkg/logger"
	"net/http"
)

//go:generate mockgen -source=adapter.go -package=api_v01 -destination=../../../mocks/adapters/http/api_v01/mock_request_builder.go
type RequestDtoBuilder interface {
	CreateCollectDto(r *http.Request) (dto.CollectDto, error)
}

type MetricAdapter struct {
	c   contracts.CollectUseCase
	g   contracts.GiveUseCase
	rb  RequestDtoBuilder
	log logger.Interface
}

func NewMetricAdapter(
	c contracts.CollectUseCase,
	g contracts.GiveUseCase,
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
