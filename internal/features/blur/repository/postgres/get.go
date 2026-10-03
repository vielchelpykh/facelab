package blur_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

func (r *BlurRepository) GetVideo(
	ctx context.Context,
	patchVideoDomain domain.VideoPatchDomain,
) (domain.VideoDomain, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, file_name, file_path, file_size, created_at
	FROM facelab.videos
	WHERE
		id=$1,
		version=$2;
	`

	row := r.Pool.QueryRow(
		ctx,
		query,
		patchVideoDomain.ID,
		patchVideoDomain.Version,
	)

	var videoModel VideoModel
	err := row.Scan(
		&videoModel.ID,
		&videoModel.Version,
		&videoModel.FileName,
		&videoModel.FilePath,
		&videoModel.FileSize,
		&videoModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.VideoDomain{}, fmt.Errorf(
				"user with id='%d': %w",
				patchVideoDomain.ID,
				err,
			)
		} else {
			return domain.VideoDomain{}, fmt.Errorf("scan error: %w", err)
		}
	}

	videoDomain := domainFromModel(videoModel)

	return videoDomain, nil
}
