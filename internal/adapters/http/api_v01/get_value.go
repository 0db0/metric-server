package api_v01

import (
	"database/sql"
	"encoding/json"
	"github.com/pkg/errors"
	"metric-server/internal/adapters/http/api_v01/dto"
	"metric-server/internal/models"
	"net/http"
	"strings"
)

func (a MetricAdapter) GetMetric(w http.ResponseWriter, r *http.Request) {
	rDto, err := a.buildDtoFromRequest(r)
	if err != nil {
		a.log.Error("Error while building metric dto", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := a.g.GetValue(r.Context(), rDto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.log.Error("Error while getting metric", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := a.buildResponseDto(metric)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		a.log.Error("Error while encoding response", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a MetricAdapter) buildDtoFromRequest(r *http.Request) (dto.ValueDto, error) {
	if r.Method == http.MethodGet {
		paths := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")
		if len(paths) != 4 {
			return dto.ValueDto{}, errors.New("invalid path")
		}

		mType := paths[len(paths)-2]
		name := paths[len(paths)-1]

		if mType == "" || name == "" {
			return dto.NewValueDto("", ""), errors.New("metric type or metric name must be non empty")
		}

		return dto.NewValueDto(name, mType), nil
	}

	var valueDto dto.ValueDto
	err := json.NewDecoder(r.Body).Decode(&valueDto)
	if err != nil {
		return dto.NewValueDto("", ""), err
	}

	return valueDto, nil
}

func (a MetricAdapter) buildResponseDto(metric models.Metric) dto.MetricResponseDto {
	return dto.MetricResponseDto{
		Name:  metric.Name,
		Type:  metric.Type,
		Delta: metric.GetDelta(),
		Value: metric.GetValue(),
	}
}
