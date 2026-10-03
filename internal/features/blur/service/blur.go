package blur_service

import (
	"context"
	"fmt"

	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

func (s *BlurService) PatchVideo(
	ctx context.Context,
	patchVideoDomain domain.VideoPatchDomain,
) (domain.VideoBlurDomain, error) {
	videoDomain, err := s.blurRepository.GetVideo(ctx, patchVideoDomain)
	if err != nil {
		return domain.VideoBlurDomain{}, fmt.Errorf("get original video from repository: %w", err)
	}

	clientDomain, err := s.blurHTTPClient.BlurVideo(ctx, videoDomain.FileName, videoDomain.FilePath)
	if err != nil {
		return domain.VideoBlurDomain{}, fmt.Errorf("get client domain: %w", err)
	}

	videoBlurDomain := domain.NewVideoBlurUninitialized(
		clientDomain.FileName,
		clientDomain.FilePath,
		clientDomain.FileSize,
		videoDomain.ID,
	)

	videoBlurDomain, err = s.blurRepository.AddBlur(ctx, videoBlurDomain)
	if err != nil {
		return domain.VideoBlurDomain{}, fmt.Errorf("add blurred video to repository: %w", err)
	}

	return videoBlurDomain, nil
}
