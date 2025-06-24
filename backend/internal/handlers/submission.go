package handlers

import (
	"net/http"

	"codegrader-backend/internal/models"
	"codegrader-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SubmissionHandler struct {
	submissionSvc services.SubmissionService
}

func NewSubmissionHandler(submissionSvc services.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{submissionSvc: submissionSvc}
}

func (h *SubmissionHandler) CreateSubmission(c *fiber.Ctx) error {
	var req models.SubmissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.FileName == "" || req.FileType == "" || req.Content == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	resp, err := h.submissionSvc.CreateSubmission(&req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(resp)
}

func (h *SubmissionHandler) GetSubmissions(c *fiber.Ctx) error {
	submissions, err := h.submissionSvc.GetAllSubmissions()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch submissions",
		})
	}

	return c.JSON(fiber.Map{
		"data": submissions,
	})
}

func (h *SubmissionHandler) GetSubmission(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing submission ID",
		})
	}

	submission, err := h.submissionSvc.GetSubmission(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Submission not found",
		})
	}

	return c.JSON(submission)
}

func (h *SubmissionHandler) DeleteSubmission(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing submission ID",
		})
	}

	if err := h.submissionSvc.DeleteSubmission(id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete submission",
		})
	}

	return c.Status(http.StatusNoContent).Send(nil)
}

func (h *SubmissionHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Service is healthy",
	})
}
