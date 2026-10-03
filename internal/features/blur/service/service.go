package blur_service

import (
	"context"

	core_http_client "github.com/vielchelpykh/facelab/internal/core/client/http"
	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

type BlurService struct {
	blurRepository BlurRepository
	blurHTTPClient *core_http_client.Client
}

type BlurRepository interface {
	AddVideo(
		video domain.VideoDomain,
		ctx context.Context,
	) (domain.VideoDomain, error)
	GetVideo(
		ctx context.Context,
		patchVideoDomain domain.VideoPatchDomain,
	) (domain.VideoDomain, error)
	AddBlur(
		ctx context.Context,
		blur domain.VideoBlurDomain,
	) (domain.VideoBlurDomain, error)
	GetBlurredByID(
		ctx context.Context,
		id int,
	) (domain.VideoBlurDomain, error)
}

func NewBlurService(
	blurRepository BlurRepository,
	blurHTTPClient *core_http_client.Client,
) *BlurService {
	return &BlurService{
		blurRepository: blurRepository,
		blurHTTPClient: blurHTTPClient,
	}
}
