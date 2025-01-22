package collect_mappers

import (
	"context"
	"metric-server/internal/adapters/http/api_v01"
	"metric-server/internal/models"
)

type GaugeMapper struct {
}

func NewGaugeMapper() GaugeMapper {
	return GaugeMapper{}
}

func (g GaugeMapper) ToMetricModel(ctx context.Context, dto api_v01.CollectDto) (models.Metric, error) {
	return models.NewMetric(dto.ID, dto.MType, dto.Delta, dto.Value), nil
}

func (g GaugeMapper) ToMetricModelList(ctx context.Context, dtoList []api_v01.CollectDto) ([]models.Metric, error) {
	var metrics []models.Metric

	for _, dto := range dtoList {
		if dto.MType == models.TypeGauge {
			metrics = append(metrics, models.NewMetric(dto.ID, dto.MType, dto.Delta, dto.Value))
		}
	}

	return metrics, nil
}
