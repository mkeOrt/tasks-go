package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mkeOrt/tasks-go/internal/config"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(cfg *config.Config, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Server.Addr,
			Handler:      handler,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	serverError := make(chan error, 1)
	go func() {
		s.logger.Info("server is starting", "addr", s.httpServer.Addr)
		serverError <- s.httpServer.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server error", "error", err)
			return err
		}
	case sig := <-quit:
		s.logger.Info("server is shutting down", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error("server shutdown error", "error", err)
			if err := s.httpServer.Close(); err != nil {
				s.logger.Error("server close error", "error", err)
				return err
			}
			return err
		}
		s.logger.Info("server is closed")
	}

	return nil
}
