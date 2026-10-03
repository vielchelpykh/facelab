package domain

type ClientDomain struct {
	FileName string `json:"file_name" validate:"required"`
	FilePath string `json:"file_path" validate:"required"`
	FileSize int64  `json:"file_size" validate:"required,gt=0"`
}
