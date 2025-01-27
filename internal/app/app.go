package app

import (
	"fmt"
	"github.com/0db0/metric-server/config"
	"github.com/0db0/metric-server/internal/adapters/db/postgres"
	"github.com/0db0/metric-server/internal/adapters/grpc"
	"github.com/0db0/metric-server/internal/adapters/http/api_v01"
	"github.com/0db0/metric-server/internal/adapters/http/api_v01/dto_builder"
	"github.com/0db0/metric-server/internal/collect_mappers"
	"github.com/0db0/metric-server/internal/pkg/logger"
	"github.com/0db0/metric-server/internal/pkg/server"
	"github.com/0db0/metric-server/internal/router"
	"github.com/0db0/metric-server/internal/usecases"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	log := logger.New()

	db := sqlx.MustOpen("pgx", cfg.DB.Dsn)
	storage := postgres.New(db)
	chainedMapper := collect_mappers.NewCounterMapper(collect_mappers.NewGaugeMapper(), storage)

	c := usecases.NewCollectUseCase(storage, chainedMapper)
	g := usecases.NewGiveUseCase(storage)
	rb := dto_builder.NewCounterBuilder(dto_builder.NewGaugeBuilder())
	v01 := api_v01.NewMetricAdapter(c, g, rb, log)

	r := router.New(v01)
	httpServer := server.NewHTTPServer(r, cfg)
	httpServer.Start()
	log.Info("start HTTP server")

	grpcAdapter := grpc.New(c, g, log)
	gRPCServer := server.NewGRPCServer(grpcAdapter, *cfg)
	gRPCServer.Start()
	log.Info("start gRPC server")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		log.Info(fmt.Sprintf("app interrupt by sygnal %s", sig.String()))
	case err := <-httpServer.Notify():
		log.Error(fmt.Errorf("http server stopped by %w", err))
	case err := <-gRPCServer.Notify():
		log.Error(fmt.Errorf("grpc server stopped by %w", err))
	}

	err := httpServer.Shutdown()

	if err != nil {
		log.Error("error occurs while shutdown: %w", err)
	}

	gRPCServer.Shutdown()

	exitCode := 0
	defer func() {
		if err := recover(); err != nil {
			log.Error("app shutdown due to panic", err)

			exitCode = 1
		}

		os.Exit(exitCode)
	}()

	log.Info("app shutdown")
}
