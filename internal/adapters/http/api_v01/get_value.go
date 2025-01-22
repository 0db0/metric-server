package api_v01

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"metric-server/internal/models"
	"net/http"
)

func (a MetricAdapter) GetMetric(w http.ResponseWriter, r *http.Request) {
	dto, err := a.buildDtoFromRequest(r)
	if err != nil {
		a.log.Error("Error while building metric dto", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := a.g.GetValue(r.Context(), dto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
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

func (a MetricAdapter) buildDtoFromRequest(r *http.Request) (ValueDto, error) {
	if r.Method == http.MethodGet {
		mType := chi.URLParam(r, "metric-type")
		name := chi.URLParam(r, "metric-name")

		if mType == "" || name == "" {
			return NewValueDto("", ""), errors.New("metric type or metric name must be non empty")
		}

		return NewValueDto(name, mType), nil
	}

	var dto ValueDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return NewValueDto("", ""), err
	}
	a.log.Debug("Got dto", dto)
	return dto, nil
}

func (a MetricAdapter) buildResponseDto(metric models.Metric) MetricResponseDto {
	return MetricResponseDto{
		Name:  metric.Name,
		Type:  metric.Type,
		Delta: metric.GetDelta(),
		Value: metric.GetValue(),
	}
}
