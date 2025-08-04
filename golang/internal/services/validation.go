package services

import (
	"sample-game-backend/internal/database"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	keystoreService = NewKeystoreService()
)

// ValidateIntent intent validation
func ValidateIntent(intent struct {
	Method string `json:"method" binding:"required"`
	From   []struct {
		Type   string `json:"type" binding:"required"`
		ID     string `json:"id" binding:"required"`
		Amount int    `json:"amount" binding:"required"`
	} `json:"from" binding:"required"`
	To []struct {
		Type   string `json:"type" binding:"required"`
		ID     string `json:"id" binding:"required"`
		Amount int    `json:"amount" binding:"required"`
	} `json:"to" binding:"required"`
}) bool {
	// Validate allowed methods
	allowedMethods := map[string]bool{
		"mint":        true,
		"transfer":    true,
		"burn":        true,
		"burn-permit": true,
	}

	if !allowedMethods[intent.Method] {
		return false
	}

	// Special validation for mint method
	if intent.Method == "mint" {
		// Must have at least one from item
		if len(intent.From) == 0 {
			return false
		}

		// All from items must be asset type
		for _, from := range intent.From {
			if from.Type != "asset" {
				return false
			}
		}
	}

	return true
}

// GenerateValidatorSignature generate validator signature (sample implementation)
func GenerateValidatorSignature(userSig hexutil.Bytes, digest common.Hash) (hexutil.Bytes, error) {
	signature, err := keystoreService.Sign(digest.Bytes())
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// ValidateAndProcessMint mint validation and processing
func ValidateAndProcessMint(sessionID string, fromAssets []struct {
	Type   string `json:"type" binding:"required"`
	ID     string `json:"id" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}) error {
	// Asset balance validation and deduction
	return database.CheckAndDeductAssets(sessionID, fromAssets)
}
