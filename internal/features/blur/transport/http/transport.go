package blur_transport_http

import (
	"context"

	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

type BlurHTTPHandler struct {
	blurService BlurService
}

type BlurService interface {
	AddVideo(
		video domain.VideoDomain,
		ctx context.Context,
	) (domain.VideoDomain, error)
	PatchVideo(
		ctx context.Context,
		patchVideoDomain domain.VideoPatchDomain,
	) (domain.VideoBlurDomain, error)
	GetBlurredByID(
		ctx context.Context,
		id int,
	) (domain.VideoBlurDomain, error)
}

func NewBlurHTTPHandler(blurService BlurService) *BlurHTTPHandler {
	return &BlurHTTPHandler{
		blurService: blurService,
	}
}
