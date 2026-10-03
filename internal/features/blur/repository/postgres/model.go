package blur_postgres_repository

import (
	"time"

	domain "github.com/vielchelpykh/facelab/internal/core/domains"
)

type VideoModel struct {
	ID        int
	Version   int
	FileName  string
	FilePath  string
	FileSize  int64
	CreatedAt time.Time
}

func domainFromModel(videoModel VideoModel) domain.VideoDomain {
	return domain.VideoDomain{
		ID:        videoModel.ID,
		Version:   videoModel.Version,
		FileName:  videoModel.FileName,
		FilePath:  videoModel.FilePath,
		FileSize:  videoModel.FileSize,
		CreatedAt: videoModel.CreatedAt,
	}
}

type BlurModel struct {
	ID              int
	Version         int
	FileName        string
	FilePath        string
	FileSize        int64
	CreatedAt       time.Time
	OriginalVideoID int
}

func blurDomainFromModel(blurModel BlurModel) domain.VideoBlurDomain {
	return domain.VideoBlurDomain{
		ID:              blurModel.ID,
		Version:         blurModel.Version,
		FileName:        blurModel.FileName,
		FilePath:        blurModel.FilePath,
		FileSize:        blurModel.FileSize,
		CreatedAt:       blurModel.CreatedAt,
		OriginalVideoID: blurModel.OriginalVideoID,
	}
}
