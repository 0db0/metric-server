package main

import (
	"context"
	"github.com/0db0/metric-server/config"
	"github.com/0db0/metric-server/internal/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"os"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	l := logger.New()

	var exitCode int
	defer func() {
		if r := recover(); r != nil {
			l.Error("error on execute migrations", r)

			exitCode = 1
		}

		cancel()
		os.Exit(exitCode)
	}()

	cfg := config.MustLoad()

	goose.SetLogger(goose.NopLogger())
	goose.SetTableName("migrations")
	err := goose.SetDialect("postgres")

	if err != nil {
		panic(errors.Wrap(err, "error on setting dialect"))
	}

	db := sqlx.MustOpen("pgx", cfg.DB.Dsn)
	err = db.Ping()
	if err != nil {
		panic(errors.Wrap(err, "error on ping to database"))
	}
	l.Info("migrations started")
	err = goose.Up(db.DB, "internal/migrations")
	if err != nil {
		panic(errors.Wrap(err, "error on running migrations"))
	}

	l.Info("migrations finished")
}
