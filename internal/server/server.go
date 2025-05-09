package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/fcjr/db/internal/store"
)

const defaultReadTimeout = 10 * time.Second
const defaultWriteTimeout = 10 * time.Second
const defaultShutdownGracePeriod = 10 * time.Second

type ServerOption func(*Server) error

type Server struct {
	logger *slog.Logger
	mux    *http.ServeMux

	readTimeout         time.Duration
	writeTimeout        time.Duration
	shutdownGracePeriod time.Duration

	store *store.Store
}

func New(opts ...ServerOption) (*Server, error) {
	s := &Server{
		logger: slog.Default(),
		mux:    http.NewServeMux(),

		readTimeout:         defaultReadTimeout,
		writeTimeout:        defaultReadTimeout,
		shutdownGracePeriod: defaultShutdownGracePeriod,

		store: store.New(),
	}

	// apply options
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	// register routes
	s.mux.HandleFunc("GET /get", s.handleGet)
	s.mux.HandleFunc("POST /set", s.handleSet)

	return s, nil
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) error {

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	s.logger.Info("server listening",
		"addr", addr,
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("server shutting down")
		// note, passing a new context here to "reset the clock" and give the grace period
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownGracePeriod)
		defer cancel()
		return httpServer.Shutdown(ctx)
	case err := <-errCh:
		s.logger.Error("server listen error, shuting down",
			"err", err,
		)
		return err
	}
}
