package handlers

import (
	"net/http"

	"sample-game-backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configure router
func SetupRoutes(r *gin.Engine) {
	// API routes configuration
	api := r.Group("/api")
	{
		// Endpoints requiring authentication
		assets := api.Group("/assets")
		assets.Use(middleware.AuthMiddleware())
		{
			assets.GET("", GetAssetsHandler)
		}

		// Order validation endpoints
		validate := api.Group("/validate")
		validate.Use(middleware.AuthMiddleware())
		validate.Use(HMACMiddleware())
		{
			validate.POST("", ValidateUserActionHandler)
		}

		result := api.Group("/result")
		result.Use(cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Authorization", "X-Dapp-Authorization", "X-Dapp-SessionID", "Content-Type", "ORIGIN", "Content-Length", "Content-Type", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With", "expires", "X-HMAC-Signature"},
		}))
		result.Use(HMACMiddleware())
		{
			result.POST("", ExchangeResultHandler)
		}
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Server is running normally",
		})
	})
}
