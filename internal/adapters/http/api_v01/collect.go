package api_v01

import (
	"encoding/json"
	"metric-server/internal/adapters/http/api_v01/dto"
	"net/http"
)

// Collect method for insert or update metrics.
// This endpoint is used to insert or update metric values by sending a POST request with the metric ID, type, and value.
// @Summary Insert or update metric value
// @Description Inserts or updates the value of a metric specified by its type, name, and value.
// This endpoint accepts a POST request with the metric ID, type, and value as path parameters.
// Supported metric types are 'gauge' and 'counter'.
// @Accept json
// @Produce json
// @Param metrics body dto.CollectDto true "Object metric to insert or update"
// @Success 200 {object} string "Metric value inserted or updated successfully"
// @Failure 400 {string} string "Bad request. Invalid metric parameters or JSON payload"
// @Router /update [post]
func (a MetricAdapter) Collect(w http.ResponseWriter, r *http.Request) {
	metric, err := a.rb.CreateCollectDto(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := a.c.CollectOne(r.Context(), metric); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// CollectFromPath method for insert or update metrics.
// This endpoint is used to insert or update metric values by sending a POST request with the metric ID, type, and value.
// @Summary Insert or update metric value
// @Description Inserts or updates the value of a metric specified by its type, name, and value.
// This endpoint accepts a POST request with the metric ID, type, and value as path parameters.
// Supported metric types are 'gauge' and 'counter'.
// @Param metricType path string true "Type of the metric ('gauge' or 'counter')"
// @Param metricName path string true "Name of the metric"
// @Param metricValue path number true "Value of the metric"
// @Produce json
// @Success 200 {string} string "Metric value inserted or updated successfully"
// @Failure 400 {string} string "Bad request. Invalid metric parameters or JSON payload"
// @Router /update/{metricType}/{metricName}/{metricValue} [post]
func (a MetricAdapter) CollectFromPath(w http.ResponseWriter, r *http.Request) {
	a.Collect(w, r)
}

// CollectMany method for bulk insert or update of metrics.
// This endpoint is used to bulk insert or update metric values by sending a POST request with a JSON array of metrics.
// @Summary Bulk insert or update metrics
// @Description Bulk inserts or updates metric values.
// This endpoint accepts a POST request with a JSON array of metrics.
// Each metric should have an ID, type, and either delta (for counter type) or value (for gauge type).
// Supported metric types are 'gauge' and 'counter'.
// @Accept json
// @Produce json
// @Param metrics body []dto.CollectDto true "Array of metrics to insert or update"
// @Success 200 {string} string "Metrics inserted or updated successfully"
// @Failure 400 {string} string "Bad request. Invalid JSON payload or metric parameters"
// @Router /updates [post]
func (a MetricAdapter) CollectMany(w http.ResponseWriter, r *http.Request) {
	var metrics []dto.CollectDto
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := a.c.CollectMany(r.Context(), metrics); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
