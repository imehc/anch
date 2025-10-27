package repository

import (
	"time"
)

// File 文件数据模型
type File struct {
	ID        int
	FileType  string // image, audio, video
	FileURL   string
	CreatedAt time.Time
}
