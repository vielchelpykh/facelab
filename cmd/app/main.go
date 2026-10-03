package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	core_http_client "github.com/vielchelpykh/facelab/internal/core/client/http"
	core_logger "github.com/vielchelpykh/facelab/internal/core/logger"
	core_postgres_pool "github.com/vielchelpykh/facelab/internal/core/repository/postgres"
	core_http_middleware "github.com/vielchelpykh/facelab/internal/core/transport/http/middleware"
	core_server_http "github.com/vielchelpykh/facelab/internal/core/transport/http/server"
	blur_postgres_repository "github.com/vielchelpykh/facelab/internal/features/blur/repository/postgres"
	blur_service "github.com/vielchelpykh/facelab/internal/features/blur/service"
	blur_transport_http "github.com/vielchelpykh/facelab/internal/features/blur/transport/http"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Логгер
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// 2. Postgres
	pool, err := core_postgres_pool.NewPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres pool", zap.Error(err))
	}
	defer pool.Close()
	logger.Info("postgres pool initialized")

	// 3. Repository
	blurRepo := blur_postgres_repository.NewBlurRepository(pool)

	// 4. HTTP client к Python
	blurServiceURL := os.Getenv("BLUR_SERVICE_URL")
	if blurServiceURL == "" {
		blurServiceURL = "http://localhost:8000"
	}
	blurHTTPClient := core_http_client.NewClient(blurServiceURL)

	// 5. Service
	blurSvc := blur_service.NewBlurService(blurRepo, blurHTTPClient)

	// 6. Handler
	blurHandler := blur_transport_http.NewBlurHTTPHandler(blurSvc)

	// 7. Роутер v1
	routerV1 := core_server_http.NewHTTPRouter(core_server_http.ApiVersion1)
	routerV1.RegisterRoutes(
		core_server_http.Route{
			Method:  http.MethodPost,
			Path:    "/videos",
			Handler: blurHandler.AddVideo,
		},
		core_server_http.Route{
			Method:  http.MethodPatch,
			Path:    "/videos",
			Handler: blurHandler.BlurVideo,
		},
		core_server_http.Route{
			Method:  http.MethodGet,
			Path:    "/videos/{id}/file",
			Handler: blurHandler.DownloadVideo,
		},
	)

	// 8. Сервер — СНАЧАЛА создаём
	cfg := core_server_http.NewConfigMust()
	httpServer := core_server_http.NewServer(cfg, logger)

	// 9. Middleware — ПОТОМ добавляем
	httpServer.WithMiddleware(
		core_http_middleware.CORS(),
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	// 10. Регистрируем роутеры
	httpServer.RegisterAPIRouters(routerV1)

	// 11. Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		logger.Warn("received signal", zap.String("signal", sig.String()))
		cancel()
	}()

	// 12. Запуск
	logger.Info("HTTP server starting", zap.String("addr", cfg.Address))
	if err := httpServer.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Fatal("HTTP server error", zap.Error(err))
	}

	logger.Info("application stopped")
}
