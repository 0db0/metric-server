package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	Config struct {
		App  App  `yaml:"app"`
		HTTP HTTP `yaml:"http"`
		DB   DB
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port            string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		ReadTimeout     time.Duration `yaml:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	}

	DB struct {
		Dsn     string `env-required:"true" env:"DATABASE_DSN"`
		PoolMax int    `yaml:"pool_max"`
	}
)

func MustLoad() *Config {
	cfg := &Config{}

	cleanenv.ReadConfig("./config/config.yaml", cfg)
	//if err := cleanenv.ReadConfig("./config/config.yaml", cfg); err != nil {
	//	panic(fmt.Errorf("yaml config error: %w", err))
	//}

	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		panic(fmt.Errorf("env config error: %w", err))
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(err)
	}

	return cfg
}
