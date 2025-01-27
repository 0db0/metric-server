package server

import (
	"context"
	"metric-server/config"
	"net/http"
	"time"
)

type HTTPServer struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewHTTPServer(h http.Handler, cfg *config.Config) *HTTPServer {
	httpServer := &http.Server{
		Handler:      h,
		Addr:         ":" + cfg.HTTP.Port,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	return &HTTPServer{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: cfg.HTTP.ShutdownTimeout,
	}
}

func (s *HTTPServer) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *HTTPServer) Notify() <-chan error {
	return s.notify
}

func (s *HTTPServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
