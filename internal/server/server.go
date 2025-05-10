package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/fcjr/db/internal/server/middleware"
	"github.com/fcjr/db/internal/store"
)

const defaultReadTimeout = 10 * time.Second
const defaultWriteTimeout = 10 * time.Second
const defaultShutdownGracePeriod = 10 * time.Second

type ServerOption func(*Server) error

type Server struct {
	logger  *slog.Logger
	handler http.Handler

	readTimeout         time.Duration
	writeTimeout        time.Duration
	shutdownGracePeriod time.Duration

	store *store.Store
}

func New(opts ...ServerOption) (*Server, error) {
	s := &Server{
		logger: slog.Default(),

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

	// create and set handler
	mux := http.NewServeMux()
	s.handler = mux

	// register routes
	mux.HandleFunc("GET /get", s.handleGet)
	mux.HandleFunc("POST /set", s.handleSet)

	// apply base middleware
	baseMiddleware := []middleware.Middleware{
		middleware.WithRecovery(s.logger), // should be outermost (first) to catch panics in other middleware
		middleware.WithLogging(s.logger),
	}
	for i := len(baseMiddleware) - 1; i >= 0; i-- {
		s.handler = baseMiddleware[i](s.handler)
	}

	return s, nil
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) error {

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      s.handler,
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
