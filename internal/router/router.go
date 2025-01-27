package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"metric-server/internal/adapters/http/api_v01"
	_ "metric-server/swagger"
)

func New(v01 *api_v01.MetricAdapter) *chi.Mux {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	// routing
	r.Route("/v0.1", func(r chi.Router) {
		// actual POST /value. GET /value/ stay for backward compatibility
		r.Get("/value/{metric-type}/{metric-name}", v01.GetMetric)
		r.Post("/value", v01.GetMetric)

		// actual POST /update. POST /update/{metric-type/{metric-name}/{value} stay for backward compatibility
		r.Post("/update/{metric-type}/{metric-name}/{value}", v01.CollectFromPath)
		r.Post("/update", v01.Collect)
		r.Post("/updates", v01.CollectMany)

		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})

	return r
}
