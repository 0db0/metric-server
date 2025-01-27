package grpc

import (
	"github.com/0db0/metric-server/internal/dto"
	pb "github.com/0db0/metric-server/pkg/metric"
)

func buildCollectDto(request *pb.CollectMetricRequest) dto.CollectDto {
	delta := request.Delta.GetValue()
	value := request.Value.GetValue()

	return dto.CollectDto{
		ID:    request.Name,
		MType: request.Type,
		Delta: &delta,
		Value: &value,
	}
}

func buildBatchCollectDto(request *pb.BatchCollectMetricRequest) []dto.CollectDto {
	var dtoList []dto.CollectDto

	for _, metric := range request.Metrics {
		dtoList = append(dtoList, buildCollectDto(metric))
	}

	return dtoList
}

func buildGetValueDto(request *pb.GetMetricRequest) dto.ValueDto {
	return dto.NewValueDto(request.Name, request.Type)
}
