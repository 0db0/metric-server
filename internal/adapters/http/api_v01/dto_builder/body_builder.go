package dto_builder

import (
	"encoding/json"
	"errors"
	"github.com/0db0/metric-server/internal/dto"
	"net/http"
	"strings"
)

type GaugeBuilder struct{}

func NewGaugeBuilder() GaugeBuilder {
	return GaugeBuilder{}
}

func (g GaugeBuilder) CreateCollectDto(r *http.Request) (dto.CollectDto, error) {
	var metric dto.CollectDto
	paths := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")

	if len(paths) == 2 {
		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			return dto.CollectDto{}, err
		}

		return metric, nil
	}

	return metric, errors.New("invalid path")
}
