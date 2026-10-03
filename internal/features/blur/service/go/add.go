package blur_service

import (
	"context"
	"fmt"

	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

func (s *BlurService) AddVideo(
	video domain.VideoDomain,
	ctx context.Context,
) (domain.VideoDomain, error) {
	video, err := s.blurRepository.AddVideo(video, ctx)
	if err != nil {
		return domain.VideoDomain{}, fmt.Errorf("add video in repository: %w", err)
	}

	return video, nil
}
