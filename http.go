package svc

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"
)

var _ Worker = (*httpServer)(nil)

// httpServer defines the internal HTTP Server worker.
type httpServer struct {
	logger     *slog.Logger
	addr       string
	httpServer *http.Server
}

func newHTTPServer(port string, handler http.Handler, logger *log.Logger) *httpServer {
	addr := net.JoinHostPort("", port)
	return &httpServer{
		addr: addr,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ErrorLog:          logger,
			ReadHeaderTimeout: 5 * time.Second, // https://medium.com/a-journey-with-go/go-understand-and-mitigate-slowloris-attack-711c1b1403f6
		},
	}
}

// Init implements the Worker interface.
func (s *httpServer) Init(logger *slog.Logger) error {
	s.logger = logger

	return nil
}

// Healthy implements the Healther interface.
func (s *httpServer) Healthy() error {
	return nil
}

// Run implements the Worker interface.
func (s *httpServer) Run() error {
	s.logger.Info("Listening and serving HTTP", slog.String("address", s.addr))
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Error("Failed to serve HTTP", slog.Any("error", err))
	}
	return nil
}

// Terminate implements the Worker interface.
func (s *httpServer) Terminate() error {
	return s.httpServer.Shutdown(context.Background())
}
