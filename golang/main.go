package main

import (
	"log/slog"

	"sample-game-backend/internal/config"
	"sample-game-backend/internal/database"
	"sample-game-backend/internal/handlers"
	"sample-game-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize database
	err := database.InitDB()
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		panic(err)
	}
	defer database.CloseDB()

	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Setup routes
	handlers.SetupRoutes(r)

	println("Server started on port 8080")
	println("API endpoint: http://localhost:8080/api/assets?language=ko")
	println("Order validation API: http://localhost:8080/api/validate")
	println("Health check: http://localhost:8080/health")
	println("Session-specific asset information is stored in go-memdb")

	r.Run(cfg.Port)
}
