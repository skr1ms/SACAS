package main

import (
	"log"

	"codegrader-backend/internal/config"
	"codegrader-backend/internal/database"
	"codegrader-backend/internal/handlers"
	"codegrader-backend/internal/repositories"
	"codegrader-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	submissionRepo := repositories.NewSubmissionRepository(db)
	openaiSvc := services.NewOpenAIService(cfg)
	submissionSvc := services.NewSubmissionService(submissionRepo, openaiSvc)
	submissionHandler := handlers.NewSubmissionHandler(submissionSvc)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	setupRoutes(app, submissionHandler)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}

func setupRoutes(app *fiber.App, submissionHandler *handlers.SubmissionHandler) {
	app.Get("/health", submissionHandler.HealthCheck)

	api := app.Group("/api")

	submissions := api.Group("/submissions")
	submissions.Post("/", submissionHandler.CreateSubmission)
	submissions.Get("/", submissionHandler.GetSubmissions)
	submissions.Get("/:id", submissionHandler.GetSubmission)
	submissions.Delete("/:id", submissionHandler.DeleteSubmission)

	api.Post("/submit", submissionHandler.CreateSubmission)
}
