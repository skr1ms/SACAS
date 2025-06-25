package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"codegrader-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSubmissionService struct {
	mock.Mock
}

func (m *MockSubmissionService) CreateSubmission(req *models.SubmissionRequest) (*models.SubmissionResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubmissionResponse), args.Error(1)
}

func (m *MockSubmissionService) GetAllSubmissions() ([]models.SubmissionListResponse, error) {
	args := m.Called()
	return args.Get(0).([]models.SubmissionListResponse), args.Error(1)
}

func (m *MockSubmissionService) GetSubmission(id string) (*models.CodeSubmission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CodeSubmission), args.Error(1)
}

func (m *MockSubmissionService) DeleteSubmission(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestSubmissionHandler_CreateSubmission_Success(t *testing.T) {
	mockService := new(MockSubmissionService)
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Post("/submissions", handler.CreateSubmission)

	reqBody := models.SubmissionRequest{
		FileName: "test.go",
		FileType: "go",
		Content:  "package main\nfunc main() {}",
	}

	expectedResponse := &models.SubmissionResponse{
		ID:       "test-id",
		Grade:    85,
		Feedback: "Good solution!",
	}

	mockService.On("CreateSubmission", &reqBody).Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/submissions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Act - выполнение тестируемого действия
	resp, err := app.Test(req)

	// Assert - проверка результатов
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response models.SubmissionResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.Grade, response.Grade)

	mockService.AssertExpectations(t)
}

func TestSubmissionHandler_CreateSubmission_InvalidBody(t *testing.T) {
	mockService := new(MockSubmissionService)
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Post("/submissions", handler.CreateSubmission)

	req := httptest.NewRequest("POST", "/submissions", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSubmissionHandler_CreateSubmission_MissingFields(t *testing.T) {
	mockService := new(MockSubmissionService)
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Post("/submissions", handler.CreateSubmission)

	reqBody := models.SubmissionRequest{
		FileName: "",
		FileType: "go",
		Content:  "package main\nfunc main() {}",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/submissions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSubmissionHandler_CreateSubmission_ServiceError(t *testing.T) {
	mockService := new(MockSubmissionService)
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Post("/submissions", handler.CreateSubmission)

	reqBody := models.SubmissionRequest{
		FileName: "test.go",
		FileType: "go",
		Content:  "package main\nfunc main() {}",
	}

	mockService.On("CreateSubmission", &reqBody).Return(nil, errors.New("database error"))

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/submissions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestSubmissionHandler_GetSubmissions_Success(t *testing.T) {
	mockService := new(MockSubmissionService)
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Get("/submissions", handler.GetSubmissions)

	expectedSubmissions := []models.SubmissionListResponse{
		{ID: "1", FileName: "test1.go", Grade: 85},
		{ID: "2", FileName: "test2.go", Grade: 90},
	}

	mockService.On("GetAllSubmissions").Return(expectedSubmissions, nil)

	req := httptest.NewRequest("GET", "/submissions", nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestSubmissionHandler_GetSubmission_Success(t *testing.T) {
	mockService := new(MockSubmissionService)
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Get("/submissions/:id", handler.GetSubmission)

	expectedSubmission := &models.CodeSubmission{
		ID:       "test-id",
		FileName: "test.go",
		Grade:    95,
	}

	mockService.On("GetSubmission", "test-id").Return(expectedSubmission, nil)

	req := httptest.NewRequest("GET", "/submissions/test-id", nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestSubmissionHandler_GetSubmission_NotFound(t *testing.T) {
	mockService := new(MockSubmissionService)
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Get("/submissions/:id", handler.GetSubmission)

	mockService.On("GetSubmission", "nonexistent").Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/submissions/nonexistent", nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestSubmissionHandler_HealthCheck(t *testing.T) {
	mockService := new(MockSubmissionService)
	handler := NewSubmissionHandler(mockService)

	app := fiber.New()
	app.Get("/health", handler.HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
