package domain

type ClientDomain struct {
	FileName string `json:"fileName" validate:"required"`
	FilePath string `json:"filePath" validate:"required"`
	FileSize int64  `json:"fileSize" validate:"required, dt=0"`
}
