package main

import (
	"metric-server/config"
	"metric-server/internal/app"
)

//	@title			Metric Collector API
//	@version		0.1
//	@description	Metrics and Alerting Service
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	Metric Collector API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	budkodmv@gmail.com
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8889
// @BasePath	/v0.1
func main() {
	cfg := config.MustLoad()
	app.Run(cfg)
}
