package blur_service

import (
	"context"
	"fmt"

	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

func (s *BlurService) GetBlurredByID(
	ctx context.Context,
	id int,
) (domain.VideoBlurDomain, error) {
	video, err := s.blurRepository.GetBlurredByID(ctx, id)
	if err != nil {
		return domain.VideoBlurDomain{}, fmt.Errorf("get blurred from repository: %w", err)
	}

	return video, nil
}
