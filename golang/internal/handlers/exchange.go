package handlers

import (
	"log/slog"
	"net/http"

	"sample-game-backend/internal/database"
	"sample-game-backend/internal/models"
	"sample-game-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// ExchangeResultHandler process result handler
func ExchangeResultHandler(c *gin.Context) {

	// Read request body
	var req models.ExchangeResp
	if err := c.ShouldBindJSON(&req); err != nil {
		LogError(slog.Default(), "ResultHandler", err, "action", "Failed to bind request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind request body"})
		return
	}

	// TODO: We need a defense logic to prevent duplicate UUIDs in requests.

	// Log request body
	LogInfo(slog.Default(), "ResultHandler", "requestBody", req)

	if len(req.Intent.Outputs) > 0 {
		// Get SessionID by UUID
		sessionID, err := database.GetSessionIDByUUID(req.UUID)
		if err != nil {
			LogError(slog.Default(), "ResultHandler", err, "action", "Failed to get session ID by UUID", "uuid", req.UUID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID or session not found"})
			return
		}

		// Process exchange result
		receiptStatus := uint64(req.Receipt.Status)
		err = services.ProcessExchangeResult(sessionID, req.Intent.Outputs, receiptStatus)
		if err != nil {
			LogError(slog.Default(), "ResultHandler", err, "action", "Failed to process exchange result")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process exchange result"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
