package blur_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

func (r *BlurRepository) GetBlurredByID(
	ctx context.Context,
	id int,
) (domain.VideoBlurDomain, error) {
	query := `
	SELECT id, version, file_name, file_path, file_size, created_at, original_video_id
	FROM facelab.blur
	WHERE id=$1;
	`
	row := r.Pool.QueryRow(ctx, query, id)

	var blurModel BlurModel
	err := row.Scan(
		&blurModel.ID,
		&blurModel.Version,
		&blurModel.FileName,
		&blurModel.FilePath,
		&blurModel.FileSize,
		&blurModel.CreatedAt,
		&blurModel.OriginalVideoID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.VideoBlurDomain{}, fmt.Errorf("not found blurred video with id: %d: %w", id, err)
		}

		return domain.VideoBlurDomain{}, fmt.Errorf("scan blurred video from database to model: %w", err)
	}

	return blurDomainFromModel(blurModel), nil
}
