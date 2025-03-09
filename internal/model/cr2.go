package model

import "time"

type CR2Backup struct {
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
	Total     int       `json:"total"`
	Processed int       `json:"processed"`
	Failed    int       `json:"failed"`
}

type CR2UploadRequest struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CR2UploadResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Filename  string    `json:"filename"`
	Filesize  int64     `json:"filesize"`
	MimeType  string    `json:"mime_type"`
	BucketURL string    `json:"bucket_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
