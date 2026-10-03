package blur_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

func (r *BlurRepository) AddVideo(
	video domain.VideoDomain,
	ctx context.Context,
) (domain.VideoDomain, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO facelab.videos (file_name, file_path, file_size, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, version, file_name, file_path, file_size, created_at
	`

	row := r.Pool.QueryRow(
		ctx,
		query,
		video.FileName,
		video.FilePath,
		video.FileSize,
		video.CreatedAt,
	)

	var videoModel VideoModel
	if err := row.Scan(
		&videoModel.ID,
		&videoModel.Version,
		&videoModel.FileName,
		&videoModel.FilePath,
		&videoModel.FileSize,
		&videoModel.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.VideoDomain{}, fmt.Errorf("not found video with id: %d: %w", videoModel.ID, err)
		}
		return domain.VideoDomain{}, fmt.Errorf("scan video from database to model: %w", err)
	}

	videoDomain := domainFromModel(videoModel)

	return videoDomain, nil
}
