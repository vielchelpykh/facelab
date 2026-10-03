package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/vielchelpykh/facelab/internal/core/errors"
	core_logger "github.com/vielchelpykh/facelab/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	rw  http.ResponseWriter
	log *core_logger.Logger
}

func NewHTTPResponseHandler(
	rw http.ResponseWriter,
	log *core_logger.Logger,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		rw:  rw,
		log: log,
	}
}

func (r *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	r.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(r.rw).Encode(responseBody); err != nil {
		r.log.Error("write HTTP response", zap.Error(err))
	}
}

func (r *HTTPResponseHandler) ErrorResponse(message string, err error) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = r.log.Warn
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = r.log.Debug
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = r.log.Warn
	default:
		statusCode = http.StatusInternalServerError
		logFunc = r.log.Error
	}

	logFunc(message, zap.Error(err))
	r.errorResponse(message, statusCode, err)
}

func (r *HTTPResponseHandler) errorResponse(
	message string,
	statusCode int,
	err error,
) {
	response := map[string]string{
		"message": message,
		"error":   err.Error(),
	}

	r.JSONResponse(response, statusCode)
}

func (r *HTTPResponseHandler) PanicResponse(message string, p any) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)
	r.log.Error(message, zap.Error(err))
	r.errorResponse(message, statusCode, err)
}
