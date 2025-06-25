package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"codegrader-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type SimpleMockService struct{}

func (m *SimpleMockService) CreateSubmission(req *models.SubmissionRequest) (*models.SubmissionResponse, error) {
	return &models.SubmissionResponse{
		ID:       "bench-id",
		Grade:    88,
		Feedback: "Benchmark test passed",
	}, nil
}

func (m *SimpleMockService) GetAllSubmissions() ([]models.SubmissionListResponse, error) {
	submissions := make([]models.SubmissionListResponse, 100)
	for i := 0; i < 100; i++ {
		submissions[i] = models.SubmissionListResponse{
			ID:       string(rune(i)),
			FileName: "test.go",
			Grade:    85,
		}
	}
	return submissions, nil
}

func (m *SimpleMockService) GetSubmission(id string) (*models.CodeSubmission, error) {
	return &models.CodeSubmission{
		ID:       id,
		FileName: "test.go",
		Grade:    85,
	}, nil
}

func (m *SimpleMockService) DeleteSubmission(id string) error {
	return nil
}

// BenchmarkCreateSubmission - бенчмарк для создания заявки
func BenchmarkCreateSubmission(b *testing.B) {
	mockService := &SimpleMockService{}
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Post("/submissions", handler.CreateSubmission)

	reqBody := models.SubmissionRequest{
		FileName: "benchmark.go",
		FileType: "go",
		Content:  "package main\nfunc main() { println(\"Hello World\") }",
	}

	jsonBody, _ := json.Marshal(reqBody)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/submissions", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// BenchmarkHealthCheck - бенчмарк для health check
func BenchmarkHealthCheck(b *testing.B) {
	mockService := &SimpleMockService{}
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Get("/health", handler.HealthCheck)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/health", nil)

		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// BenchmarkGetSubmissions - бенчмарк для получения списка заявок
func BenchmarkGetSubmissions(b *testing.B) {
	mockService := &SimpleMockService{}
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Get("/submissions", handler.GetSubmissions)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/submissions", nil)

		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}
