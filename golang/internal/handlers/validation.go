package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"sample-game-backend/internal/database"
	"sample-game-backend/internal/models"
	"sample-game-backend/internal/services"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
)

// ValidateUserActionHandler user action validation handler
func ValidateUserActionHandler(c *gin.Context) {
	var req models.ValidateRequest

	// Request binding and validation
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidateErrorResponse(c, http.StatusBadRequest, ErrorCodeInvalidRequest)
		return
	}

	// Intent validation
	if !services.ValidateIntent(req.Intent) {
		ValidateErrorResponse(c, http.StatusBadRequest, ErrorCodeInvalidIntent)
		return
	}

	// Get session ID
	sessionID, valid := ValidateSessionID(c)
	if !valid {
		return
	}

	// Store UUID and SessionID mapping
	err := database.StoreUUIDMapping(req.UUID, sessionID)
	if err != nil {
		LogError(slog.Default(), "ValidateUserActionHandler", err, "action", "Failed to store UUID mapping")
		ValidateErrorResponse(c, http.StatusInternalServerError, ErrorCodeUUIDMappingFailed)
		return
	}

	requestBytes, _ := json.Marshal(req)
	LogInfo(slog.Default(), "ValidateUserActionHandler", "sessionID", sessionID, "uuid", req.UUID, "req", string(requestBytes))

	// For mint method, validate and deduct assets
	if req.Intent.Type == "assemble" {
		if err := services.ValidateAndProcessMint(sessionID, req.Intent.From); err != nil {
			ValidateErrorResponse(c, http.StatusBadRequest, ErrorCodeInsufficientBalance)
			return
		}
	}

	// Generate validator signature (in actual implementation, use validator's private key)
	userSigBytes := hexutil.MustDecode(req.UserSig)
	digestHash := common.HexToHash(req.Digest)
	validatorSig, err := services.GenerateValidatorSignature(userSigBytes, digestHash)
	if err != nil {
		LogError(slog.Default(), "GenerateValidatorSignature", err)
		ValidateErrorResponse(c, http.StatusInternalServerError, ErrorCodeSignatureGeneration)
		return
	}

	LogInfo(slog.Default(), "validateUserActionHandler", "validatorSig", validatorSig, "userSig", req.UserSig, "digest", req.Digest)

	// Success response
	response := models.ValidateResponse{
		Success: true,
		Data: struct {
			UserSig      string `json:"userSig"`
			ValidatorSig string `json:"validatorSig"`
		}{
			UserSig:      req.UserSig,
			ValidatorSig: validatorSig.String(),
		},
	}

	c.JSON(http.StatusOK, response)
}
