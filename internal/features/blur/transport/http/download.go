package blur_transport_http

import (
	"net/http"
	"strconv"

	core_errors "github.com/vielchelpykh/facelab/internal/core/errors"
	core_logger "github.com/vielchelpykh/facelab/internal/core/logger"
	core_http_response "github.com/vielchelpykh/facelab/internal/core/response"
)

func (h *BlurHTTPHandler) DownloadVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(w, log)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseHandler.ErrorResponse("get user id", core_errors.ErrInvalidArgument)
		return
	}

	video, err := h.blurService.GetBlurredByID(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse("get video from service", err)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeFile(w, r, video.FilePath)
}
