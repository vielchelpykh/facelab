package blur_transport_http

import (
	"net/http"
	"time"

	domain "github.com/vielchelpykh/facelab/internal/core/domains"
	core_errors "github.com/vielchelpykh/facelab/internal/core/errors"
	core_logger "github.com/vielchelpykh/facelab/internal/core/logger"
	core_http_request "github.com/vielchelpykh/facelab/internal/core/request"
	core_http_response "github.com/vielchelpykh/facelab/internal/core/response"
)

type PatchVideoDTO struct {
	ID      int `json:"id" validate:"required,gt=0"`
	Version int `json:"version" validate:"required,gt=0"`
}

func patchDomainFromDTO(patchVideoDTO PatchVideoDTO) domain.VideoPatchDomain {
	return domain.NewVideoPatchDomain(
		patchVideoDTO.ID,
		patchVideoDTO.Version,
	)
}

type BlurVideoDTO struct {
	ID              int       `json:"id"`
	Version         int       `json:"version"`
	FileName        string    `json:"file_name"`
	FilePath        string    `json:"file_path"`
	FileSize        int64     `json:"file_size"`
	CreatedAt       time.Time `json:"created_at"`
	OriginalVideoID int       `json:"original_video_id"`
}

func dtoBlurFromDomain(videoBlurDomain domain.VideoBlurDomain) BlurVideoDTO {
	return BlurVideoDTO{
		ID:              videoBlurDomain.ID,
		Version:         videoBlurDomain.Version,
		FileName:        videoBlurDomain.FileName,
		FilePath:        videoBlurDomain.FilePath,
		FileSize:        videoBlurDomain.FileSize,
		CreatedAt:       videoBlurDomain.CreatedAt,
		OriginalVideoID: videoBlurDomain.OriginalVideoID,
	}
}

func (h *BlurHTTPHandler) BlurVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(w, log)

	var patchVideoDTO PatchVideoDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &patchVideoDTO); err != nil {
		responseHandler.ErrorResponse("decode and validate request", core_errors.ErrInvalidArgument)
		return
	}

	patchVideoDomain := patchDomainFromDTO(patchVideoDTO)

	videoBlurDomain, err := h.blurService.PatchVideo(ctx, patchVideoDomain)
	if err != nil {
		responseHandler.ErrorResponse("get video blur domain from service", err)
		return
	}

	response := dtoBlurFromDomain(videoBlurDomain)

	responseHandler.JSONResponse(response, http.StatusOK)
}
