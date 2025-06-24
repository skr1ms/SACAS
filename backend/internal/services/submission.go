package services

import (
	"fmt"
	"log"
	"time"

	"codegrader-backend/internal/models"
	"codegrader-backend/internal/repositories"

	"github.com/google/uuid"
)

type SubmissionService interface {
	CreateSubmission(req *models.SubmissionRequest) (*models.SubmissionResponse, error)
	GetSubmission(id string) (*models.CodeSubmission, error)
	GetAllSubmissions() ([]models.SubmissionListResponse, error)
	DeleteSubmission(id string) error
}

type submissionService struct {
	repo      repositories.SubmissionRepository
	openaiSvc OpenAIService
}

func NewSubmissionService(repo repositories.SubmissionRepository, openaiSvc OpenAIService) SubmissionService {
	return &submissionService{
		repo:      repo,
		openaiSvc: openaiSvc,
	}
}

func (s *submissionService) CreateSubmission(req *models.SubmissionRequest) (*models.SubmissionResponse, error) {
	allowedTypes := []string{".cpp", ".java", ".js", ".kt", ".py"}
	if !contains(allowedTypes, req.FileType) {
		return nil, fmt.Errorf("unsupported file type: %s", req.FileType)
	}

	submission := &models.CodeSubmission{
		ID:        uuid.New().String(),
		FileName:  req.FileName,
		FileType:  req.FileType,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	grade, feedback, err := s.openaiSvc.AnalyzeCode(req.Content, req.FileType)
	if err != nil {
		log.Printf("OpenAI analysis failed: %v", err)
		grade = 3
		feedback = "Автоматический анализ недоступен. Код загружен для ручной проверки."
	}

	submission.Grade = grade
	submission.Feedback = feedback

	if err := s.repo.Create(submission); err != nil {
		return nil, err
	}

	return &models.SubmissionResponse{
		ID:       submission.ID,
		Grade:    submission.Grade,
		Feedback: submission.Feedback,
	}, nil
}

func (s *submissionService) GetSubmission(id string) (*models.CodeSubmission, error) {
	return s.repo.GetByID(id)
}

func (s *submissionService) GetAllSubmissions() ([]models.SubmissionListResponse, error) {
	submissions, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]models.SubmissionListResponse, len(submissions))
	for i, sub := range submissions {
		result[i] = models.SubmissionListResponse{
			ID:        sub.ID,
			FileName:  sub.FileName,
			FileType:  sub.FileType,
			Grade:     sub.Grade,
			CreatedAt: sub.CreatedAt,
		}
	}

	return result, nil
}

func (s *submissionService) DeleteSubmission(id string) error {
	return s.repo.Delete(id)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
