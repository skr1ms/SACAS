package models

import (
	"time"
)

type CodeSubmission struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	FileName  string    `json:"file_name" gorm:"not null"`
	FileType  string    `json:"file_type" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Grade     int       `json:"grade"`
	Feedback  string    `json:"feedback" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type SubmissionRequest struct {
	FileName string `json:"file_name" validate:"required"`
	FileType string `json:"file_type" validate:"required,oneof=.cpp .java .js .kt .py"`
	Content  string `json:"content" validate:"required"`
}

type SubmissionResponse struct {
	ID       string `json:"id"`
	Grade    int    `json:"grade"`
	Feedback string `json:"feedback"`
}

type SubmissionListResponse struct {
	ID        string    `json:"id"`
	FileName  string    `json:"file_name"`
	FileType  string    `json:"file_type"`
	Grade     int       `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
}
