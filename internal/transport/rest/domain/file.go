package domain

import (
	"time"
)

type (
	FileStatus int
	FileType   string
)

const (
	ClientUploadInProgress FileStatus = iota
	UploadedByClient
	ClientUploadError
	StorageUploadInProgress
	UploadedToStorage
	StorageUploadError
)

const (
	Image FileType = "image"
)

type File struct {
	ID              string     `json:"id"`
	ProductId       string     `json:"product_id"`
	Type            FileType   `json:"type"`
	ContentType     string     `json:"content_type"`
	Name            string     `json:"name"`
	Size            int64      `json:"size"`
	Status          FileStatus `json:"status"`
	UploadStartedAt time.Time  `json:"upload_started_at"`
	URL             string     `json:"url"`
}
