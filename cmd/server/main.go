package main

import (
	"metric-server/config"
	"metric-server/internal/app"
)

func main() {
	cfg := config.MustLoad()
	app.Run(cfg)
}
