package repositories

import (
	"codegrader-backend/internal/models"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	Create(submission *models.CodeSubmission) error
	GetByID(id string) (*models.CodeSubmission, error)
	GetAll() ([]models.CodeSubmission, error)
	Update(submission *models.CodeSubmission) error
	Delete(id string) error
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Create(submission *models.CodeSubmission) error {
	return r.db.Create(submission).Error
}

func (r *submissionRepository) GetByID(id string) (*models.CodeSubmission, error) {
	var submission models.CodeSubmission
	err := r.db.First(&submission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *submissionRepository) GetAll() ([]models.CodeSubmission, error) {
	var submissions []models.CodeSubmission
	err := r.db.Order("created_at DESC").Find(&submissions).Error
	return submissions, err
}

func (r *submissionRepository) Update(submission *models.CodeSubmission) error {
	return r.db.Save(submission).Error
}

func (r *submissionRepository) Delete(id string) error {
	return r.db.Delete(&models.CodeSubmission{}, "id = ?", id).Error
}
