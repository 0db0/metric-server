package api_v01

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (a MetricAdapter) Collect(w http.ResponseWriter, r *http.Request) {
	var metric CollectDto
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := a.c.CollectOne(r.Context(), metric); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (a MetricAdapter) CollectMany(w http.ResponseWriter, r *http.Request) {
	var metrics []CollectDto
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := a.c.CollectMany(r.Context(), metrics); err != nil {
		a.log.Debug("DEBUD COLLECT MANY", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (a MetricAdapter) buildCollectDto(r *http.Request) (CollectDto, error) {
	var metric CollectDto

	if r.Method == http.MethodGet {
		name := chi.URLParam(r, "metric-name")
		mType := chi.URLParam(r, "metric-type")
		value, err := strconv.ParseFloat(chi.URLParam(r, "metric-value"), 64)
		if err != nil {
			return CollectDto{}, err
		}

		return CollectDto{ID: name, MType: mType, Value: &value}, nil
	}

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		return CollectDto{}, err
	}

	return metric, nil
}
