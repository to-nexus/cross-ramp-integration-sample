package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		dappAuth := c.GetHeader("X-Dapp-Authorization")
		sessionID := c.GetHeader("X-Dapp-SessionID")

		slog.Info("AuthMiddleware", "FullPath", c.FullPath(), "authHeader", authHeader, "dappAuth", dappAuth, "sessionID", sessionID)
		c.Set("Authorization", authHeader)
		c.Set("X-Dapp-Authorization", dappAuth)
		c.Set("X-Dapp-SessionID", sessionID)
		c.Next()
	}
}

// CORSMiddleware CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, X-Dapp-Authorization, X-Dapp-SessionID, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
