package handlers

import (
	"log/slog"
	"net/http"

	"sample-game-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// Common error codes
const (
	ErrorCodeInvalidRequest      = "INVALID_REQUEST"
	ErrorCodeInvalidSessionID    = "INVALID_SESSION_ID"
	ErrorCodeDBError             = "DB_ERROR"
	ErrorCodeInvalidIntent       = "INVALID_INTENT"
	ErrorCodeUUIDMappingFailed   = "UUID_MAPPING_FAILED"
	ErrorCodeInsufficientBalance = "INSUFFICIENT_BALANCE"
	ErrorCodeSignatureGeneration = "SIGNATURE_GENERATION_FAILED"
)

// ErrorResponse creates a standard error response
func ErrorResponse(c *gin.Context, statusCode int, errorCode string) {
	response := models.Response{
		Success:   false,
		ErrorCode: &errorCode,
	}
	c.JSON(statusCode, response)
}

// ValidateErrorResponse creates a standard validate error response
func ValidateErrorResponse(c *gin.Context, statusCode int, errorCode string) {
	response := models.ValidateResponse{
		Success:   false,
		ErrorCode: &errorCode,
	}
	c.JSON(statusCode, response)
}

// LogError logs an error with consistent formatting
func LogError(logger *slog.Logger, message string, err error, fields ...any) {
	logger.Error(message, append([]any{"error", err}, fields...)...)
}

// LogInfo logs info with consistent formatting
func LogInfo(logger *slog.Logger, message string, fields ...any) {
	logger.Info(message, fields...)
}

// GetSessionIDFromContext extracts session ID from gin context
func GetSessionIDFromContext(c *gin.Context) string {
	return c.GetString("X-Dapp-SessionID")
}

// ValidateSessionID validates session ID and returns error response if invalid
func ValidateSessionID(c *gin.Context) (string, bool) {
	sessionID := GetSessionIDFromContext(c)
	if sessionID == "" {
		ErrorResponse(c, http.StatusBadRequest, ErrorCodeInvalidSessionID)
		return "", false
	}
	return sessionID, true
}
