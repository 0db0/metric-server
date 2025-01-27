package grpc

import (
	"context"
	"github.com/0db0/metric-server/internal/contracts"
	"github.com/0db0/metric-server/internal/pkg/logger"
	pb "github.com/0db0/metric-server/pkg/metric"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Adapter struct {
	c   contracts.CollectUseCase
	g   contracts.GiveUseCase
	log logger.Interface

	pb.UnimplementedMetricServer
}

func New(c contracts.CollectUseCase, g contracts.GiveUseCase, log logger.Interface) Adapter {
	return Adapter{
		c:   c,
		g:   g,
		log: log,
	}
}

func (a Adapter) Collect(ctx context.Context, request *pb.CollectMetricRequest) (*pb.CollectMetricResponse, error) {
	metricDto := buildCollectDto(request)
	if err := a.c.CollectOne(ctx, metricDto); err != nil {
		a.log.Error("Error while collect metric", err)
		return nil, err
	}

	return &pb.CollectMetricResponse{}, nil
}

func (a Adapter) BatchCollect(ctx context.Context, request *pb.BatchCollectMetricRequest) (*pb.BatchCollectMetricResponse, error) {
	dto := buildBatchCollectDto(request)
	if err := a.c.CollectMany(ctx, dto); err != nil {
		a.log.Error("Error while batch collect metrics", err)
		return nil, err
	}

	return &pb.BatchCollectMetricResponse{}, nil
}

func (a Adapter) GetMetric(ctx context.Context, request *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	dto := buildGetValueDto(request)
	metric, err := a.g.GetValue(ctx, dto)

	if err != nil {
		a.log.Error("Error while getting metric", err)
		return &pb.GetMetricResponse{}, err
	}

	delta := metric.GetDelta()
	var deltaWrapped *wrapperspb.Int64Value = nil

	if delta != nil {
		deltaWrapped = wrapperspb.Int64(*delta)
	}

	value := metric.GetValue()
	var valueWrapped *wrapperspb.DoubleValue = nil

	if value != nil {
		valueWrapped = wrapperspb.Double(*value)
	}

	return &pb.GetMetricResponse{
		Name:  metric.Name,
		Type:  metric.Type,
		Delta: deltaWrapped,
		Value: valueWrapped,
	}, nil
}
