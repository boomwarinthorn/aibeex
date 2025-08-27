package main

import (
	"log"

	"fiber-hello-world/config"
	"fiber-hello-world/internal/infrastructure/database"
	"fiber-hello-world/internal/presentation/handler"
	"fiber-hello-world/internal/presentation/middleware"
	"fiber-hello-world/internal/usecase"
	"fiber-hello-world/pkg/jwt"
	"fiber-hello-world/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "fiber-hello-world/docs" // This line is needed for go-swagger to find your docs!
)

// @title Fiber Authentication API
// @version 2.0
// @description A Go Fiber API with JWT authentication, user registration, and login functionality built with Clean Architecture
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://github.com/boomwarinthorn/aibeex/blob/main/LICENSE
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := database.NewSQLiteUserRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize services
	jwtService := jwt.NewService(cfg.JWTSecret)
	validatorService := validator.NewService()

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase, jwtService, validatorService)

	// Create fiber app
	app := fiber.New(fiber.Config{
		AppName: "Fiber Authentication API v2.0",
	})

	// Swagger documentation route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// @Summary Get hello world message
	// @Description Returns a simple hello world JSON response
	// @Tags general
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router / [get]
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World - Clean Architecture",
			"version": "2.0",
		})
	})

	// Public routes
	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)

	// Protected routes
	protected := app.Group("/", middleware.JWTMiddleware(jwtService))
	protected.Get("/me", userHandler.GetMe)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
