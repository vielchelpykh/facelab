package blur_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

func (r *BlurRepository) AddBlur(
	ctx context.Context,
	blur domain.VideoBlurDomain,
) (domain.VideoBlurDomain, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO facelab.blur (file_name, file_path, file_size, created_at, original_video_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, version, file_name, file_path, file_size, created_at, original_video_id
	`
	row := r.Pool.QueryRow(
		ctx,
		query,
		blur.FileName,
		blur.FilePath,
		blur.FileSize,
		blur.CreatedAt,
		blur.OriginalVideoID,
	)

	var blurModel BlurModel
	if err := row.Scan(
		&blurModel.ID,
		&blurModel.Version,
		&blurModel.FileName,
		&blurModel.FilePath,
		&blurModel.FileSize,
		&blurModel.CreatedAt,
		&blurModel.OriginalVideoID,
	); err != nil {
		return domain.VideoBlurDomain{}, fmt.Errorf("scan error: %w", err)
	}

	blurDomain := blurDomainFromModel(blurModel)

	return blurDomain, nil
}
