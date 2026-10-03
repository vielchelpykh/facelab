package blur_transport_http

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	domain "github.com/vielchelpykh/facelab/internal/core/domains"
	core_errors "github.com/vielchelpykh/facelab/internal/core/errors"
	core_logger "github.com/vielchelpykh/facelab/internal/core/logger"
	core_http_response "github.com/vielchelpykh/facelab/internal/core/response"
)

type OriginalVideoDTO struct {
	ID        int       `json:"id"`
	Version   int       `json:"version"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}

func dtoFromDomain(video domain.VideoDomain) (OriginalVideoDTO, error) {
	return OriginalVideoDTO{
		ID:        video.ID,
		Version:   video.Version,
		FileName:  video.FileName,
		FilePath:  video.FilePath,
		FileSize:  video.FileSize,
		CreatedAt: video.CreatedAt,
	}, nil
}

func domainFromDTO(fileName string, filePath string, fileSize int64) domain.VideoDomain {
	return domain.NewVideoUninitialized(fileName, filePath, fileSize)
}

func (h *BlurHTTPHandler) AddVideo(w http.ResponseWriter, r *http.Request) {
	var (
		folder = os.Getenv("ORIGINAL_VIDEOS_FOLDER")
	)

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(w, log)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		responseHandler.ErrorResponse("parse multipart form", core_errors.ErrInvalidArgument)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		responseHandler.ErrorResponse("get file", err)
		return
	}
	defer file.Close()

	if err := os.MkdirAll(folder, 0755); err != nil {
		responseHandler.ErrorResponse("create directory", err)
		return
	}

	filePath := filepath.Join(
		folder,
		header.Filename,
	)

	destinationVideo, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		responseHandler.ErrorResponse("open new file", err)
		return
	}
	defer destinationVideo.Close()

	if _, err := io.Copy(destinationVideo, file); err != nil {
		responseHandler.ErrorResponse("copy file", err)
		return
	}

	info, err := destinationVideo.Stat()
	if err != nil {
		responseHandler.ErrorResponse("get video information", err)
		return
	}

	videoDomain := domainFromDTO(info.Name(), filePath, info.Size())

	videoDomain, err = h.blurService.AddVideo(videoDomain, ctx)
	if err != nil {
		responseHandler.ErrorResponse("failed to add video", err)
		return
	}

	response, err := dtoFromDomain(videoDomain)
	if err != nil {
		responseHandler.ErrorResponse("failed to create response", err)
		return
	}

	responseHandler.JSONResponse(response, http.StatusCreated)
}
