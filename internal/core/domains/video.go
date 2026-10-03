package domain

import (
	"time"
)

type VideoDomain struct {
	ID        int
	Version   int
	FileName  string
	FilePath  string
	FileSize  int64
	CreatedAt time.Time
}

func NewVideo(
	id int,
	version int,
	fileName string,
	filePath string,
	fileSize int64,
) VideoDomain {
	return VideoDomain{
		ID:        id,
		Version:   version,
		FileName:  fileName,
		FilePath:  filePath,
		FileSize:  fileSize,
		CreatedAt: time.Now(),
	}
}

func NewVideoUninitialized(fileName string, filePath string, fileSize int64) VideoDomain {
	return NewVideo(
		UninitializedID,
		UninitializedVersion,
		fileName,
		filePath,
		fileSize,
	)
}

type VideoPatchDomain struct {
	ID      int
	Version int
}

func NewVideoPatchDomain(
	id int,
	version int,
) VideoPatchDomain {
	return VideoPatchDomain{
		ID:      id,
		Version: version,
	}
}

type VideoBlurDomain struct {
	ID              int
	Version         int
	FileName        string
	FilePath        string
	FileSize        int64
	CreatedAt       time.Time
	OriginalVideoID int
}

func NewVideoBlur(
	id int,
	version int,
	fileName string,
	filePath string,
	fileSize int64,
	originalVideoID int,
) VideoBlurDomain {
	return VideoBlurDomain{
		ID:              id,
		Version:         version,
		FileName:        fileName,
		FilePath:        filePath,
		FileSize:        fileSize,
		CreatedAt:       time.Now(),
		OriginalVideoID: originalVideoID,
	}
}

func NewVideoBlurUninitialized(
	fileName string,
	filePath string,
	fileSize int64,
	originalVideoID int,
) VideoBlurDomain {
	return NewVideoBlur(
		UninitializedID,
		UninitializedVersion,
		fileName,
		filePath,
		fileSize,
		originalVideoID,
	)
}
