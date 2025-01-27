package collect_mappers

import (
	"context"
	"github.com/0db0/metric-server/internal/dto"
	"github.com/0db0/metric-server/internal/models"
	"github.com/0db0/metric-server/internal/usecases"
)

type CounterMapper struct {
	m usecases.Mapper
	s usecases.Storage
}

func NewCounterMapper(m usecases.Mapper, s usecases.Storage) CounterMapper {
	return CounterMapper{
		m: m,
		s: s,
	}
}

func (c CounterMapper) ToMetricModel(ctx context.Context, dto dto.CollectDto) (models.Metric, error) {
	if dto.MType == models.TypeCounter {
		model, err := c.s.FindByNameAndType(ctx, dto.ID, dto.MType)
		if err != nil {
			return models.Metric{}, err
		}

		metricMap := c.toMetricMap([]models.Metric{model})

		return c.toMetricModel(dto, metricMap), nil
	}

	return c.m.ToMetricModel(ctx, dto)
}

func (c CounterMapper) ToMetricModelList(ctx context.Context, dtoList []dto.CollectDto) ([]models.Metric, error) {
	filtered := c.filter(dtoList)
	var metricList []models.Metric

	if len(filtered) > 0 {
		counterMetrics, err := c.s.FindAllByType(ctx, models.TypeCounter)
		if err != nil {
			return []models.Metric{}, err
		}

		metricMap := c.toMetricMap(counterMetrics)

		for _, dto := range filtered {
			metricList = append(metricList, c.toMetricModel(dto, metricMap))
		}
	}

	add, err := c.m.ToMetricModelList(ctx, dtoList)
	if err != nil {
		return []models.Metric{}, err
	}

	return append(add, metricList...), nil
}

func (c CounterMapper) filter(dtoList []dto.CollectDto) []dto.CollectDto {
	var filtered []dto.CollectDto
	for _, dto := range dtoList {
		if dto.MType == models.TypeCounter {
			filtered = append(filtered, dto)
		}
	}

	return filtered
}

func (c CounterMapper) toMetricModel(dto dto.CollectDto, mapCounter map[string]models.Metric) models.Metric {
	var delta *int64
	m, ok := mapCounter[dto.ID]

	if ok && m.Delta.Valid && dto.Delta != nil {
		sum := *m.GetDelta() + *dto.Delta
		delta = &sum
	} else {
		delta = dto.Delta
	}

	return models.NewMetric(dto.ID, dto.MType, delta, dto.Value)
}

func (c CounterMapper) toMetricMap(metrics []models.Metric) map[string]models.Metric {
	counters := make(map[string]models.Metric)

	for _, m := range metrics {
		counters[m.Name] = m
	}

	return counters
}
