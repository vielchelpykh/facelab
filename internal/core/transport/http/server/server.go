package core_server_http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/vielchelpykh/facelab/internal/core/logger"
	core_http_middleware "github.com/vielchelpykh/facelab/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	log        *core_logger.Logger
	middleware []core_http_middleware.Middleware
}

func NewServer(
	config Config,
	log *core_logger.Logger,
) *HTTPServer {
	return &HTTPServer{
		mux:    http.NewServeMux(),
		config: config,
		log:    log,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion) + "/"

		s.mux.Handle(
			prefix,
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Address,
		Handler: mux,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr", s.config.Address))
		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("close HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
