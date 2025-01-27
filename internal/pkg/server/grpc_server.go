package server

import (
	"google.golang.org/grpc"
	"metric-server/config"
	grpcAdapter "metric-server/internal/adapters/grpc"
	"metric-server/pkg/metric"
	"net"
)

type GRPCServer struct {
	server *grpc.Server
	listen *net.Listener
	notify chan error
}

func NewGRPCServer(a grpcAdapter.Adapter, cfg config.Config) *GRPCServer {
	listen, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	server.RegisterService(&metric.Metric_ServiceDesc, a)

	return &GRPCServer{
		server: server,
		listen: &listen,
		notify: make(chan error, 1),
	}
}

func (s *GRPCServer) Start() {
	go func() {
		s.notify <- s.server.Serve(*s.listen)
		close(s.notify)
	}()
}

func (s *GRPCServer) Notify() <-chan error {
	return s.notify
}

func (s *GRPCServer) Shutdown() {
	s.server.GracefulStop()
}
