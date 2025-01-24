package dto_builder

import (
	"metric-server/internal/adapters/http/api_v01"
	"metric-server/internal/adapters/http/api_v01/dto"
	"metric-server/internal/models"
	"net/http"
	"strconv"
	"strings"
)

type CounterBuilder struct {
	i api_v01.RequestDtoBuilder
}

func NewCounterBuilder(i api_v01.RequestDtoBuilder) *CounterBuilder {
	return &CounterBuilder{
		i: i,
	}
}

func (c CounterBuilder) CreateCollectDto(r *http.Request) (dto.CollectDto, error) {
	paths := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")

	if len(paths) == 5 {
		valueString := paths[len(paths)-1]

		name := paths[len(paths)-2]
		mType := paths[len(paths)-3]

		if mType == models.TypeCounter {
			delta, err := strconv.ParseInt(valueString, 10, 64)
			if err != nil {
				return dto.CollectDto{}, err
			}

			return dto.CollectDto{ID: name, MType: mType, Delta: &delta}, nil
		} else {
			value, err := strconv.ParseFloat(valueString, 64)
			if err != nil {
				return dto.CollectDto{}, err
			}

			return dto.CollectDto{ID: name, MType: mType, Value: &value}, nil
		}
	}

	return c.i.CreateCollectDto(r)
}
