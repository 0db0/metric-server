package app

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"metric-server/config"
	"metric-server/internal/adapters/db/postgres"
	"metric-server/internal/adapters/http/api_v01"
	"metric-server/internal/adapters/http/api_v01/dto_builder"
	"metric-server/internal/collect_mappers"
	"metric-server/internal/pkg/logger"
	"metric-server/internal/pkg/server"
	"metric-server/internal/router"
	"metric-server/internal/usecases"
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
	s := server.New(r, cfg)
	s.Start()

	log.Info("start HTTP server")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info(fmt.Sprintf("app interrupt by sygnal %s", s.String()))
	case err := <-s.Notify():
		log.Error(fmt.Errorf("http server stopped by %w", err))
	}

	err := s.Shutdown()

	if err != nil {
		log.Error("error occurs while shutdown: %w", err)
	}

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
