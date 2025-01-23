package api_v01

import (
	"encoding/json"
	"metric-server/internal/adapters/http/api_v01/dto"
	"net/http"
)

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
