package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/vielchelpykh/facelab/internal/core/logger"
	core_http_response "github.com/vielchelpykh/facelab/internal/core/response"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			w.Header().Set(requestIDHeader, requestID)
			r.Header.Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("URL", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseWriter := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.Time("time", before.UTC()),
				zap.String("http_method", r.Method),
			)

			next.ServeHTTP(responseWriter, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", responseWriter.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(w, log)

			go func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse("during HTTP request got unexpected panic", p)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
