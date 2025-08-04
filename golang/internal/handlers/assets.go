package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"sample-game-backend/internal/database"
	"sample-game-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// GetAssetsHandler asset information retrieval handler
func GetAssetsHandler(c *gin.Context) {
	language := c.Query("language")

	// Validate session ID
	sessionID, valid := ValidateSessionID(c)
	if !valid {
		return
	}

	// Validate language parameter
	if language == "" {
		language = "ko" // Default value
	}

	// Get or create session-specific asset information
	sessionAssets, err := database.GetOrCreateSessionAssets(sessionID)
	if err != nil {
		LogError(slog.Default(), "GetAssetsHandler", err, "sessionID", sessionID)
		ErrorResponse(c, http.StatusInternalServerError, ErrorCodeDBError)
		return
	}

	// Convert to Asset struct
	var assets []models.Asset
	for id, balance := range sessionAssets.Assets {
		assets = append(assets, models.Asset{
			ID:      id,
			Balance: balance,
		})
	}

	v1Data := models.V1Data{
		PlayerID:      sessionID,
		Name:          fmt.Sprintf("playerName_%s", sessionID),
		WalletAddress: "0xaaaa",
		Server:        "test",
		Assets:        assets,
	}

	// Parse session information
	createdAt, _ := time.Parse(time.RFC3339, sessionAssets.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, sessionAssets.UpdatedAt)

	guide := struct {
		Authorization string `json:"Authorization"`
		DappAuth      string `json:"X-Dapp-Authorization"`
		SessionID     string `json:"X-Dapp-SessionID"`
		Message       string `json:"message"`
		SessionInfo   struct {
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"session_info"`
	}{
		Authorization: c.GetString("Authorization"),
		DappAuth:      c.GetString("X-Dapp-Authorization"),
		SessionID:     sessionID,
		Message:       "The guide field displays header information at request time. It is used to verify that the game company and protocol are correctly matched and is not provided to the game company. For ramp frontend developer reference.",
		SessionInfo: struct {
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}

	response := models.Response{
		Success:   true,
		ErrorCode: nil,
		Data: struct {
			V1    models.V1Data `json:"v1"`
			Guide any           `json:"guide"`
		}{
			V1:    v1Data,
			Guide: guide,
		},
	}

	c.JSON(http.StatusOK, response)
}
