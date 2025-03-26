package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/zetsub0/yakvs/internal/config"
)

// Server implements http server
type Server struct {
	server     http.Server
	ctxTimeout time.Duration
}

// New ...
func New(cfg config.HTTPServer, handler http.Handler) *Server {

	return &Server{
		server: http.Server{
			Addr:         cfg.Address,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},

		ctxTimeout: cfg.ContextTimeout,
	}
}

// Run starts server listening and working
func (s *Server) Run(ctx context.Context) {

	go s.server.ListenAndServe()
	slog.Info("Server started", "addr", s.server.Addr)
	<-ctx.Done()

	sCTX, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()

	s.server.Shutdown(sCTX)
}
