package services

import (
	"log/slog"

	"sample-game-backend/internal/database"
	"sample-game-backend/internal/models"
)

// ProcessExchangeResult process exchange result
func ProcessExchangeResult(sessionID string, outputs []models.PairAsset, receiptStatus uint64) error {
	// Skip processing if receipt status is not 0x1
	if receiptStatus != 1 {
		slog.Info("ProcessExchangeResult", "sessionID", sessionID, "receiptStatus", receiptStatus, "action", "skipped")
		return nil
	}

	// Skip processing if output is empty
	if len(outputs) == 0 {
		slog.Info("ProcessExchangeResult", "sessionID", sessionID, "outputs", "empty", "action", "skipped")
		return nil
	}

	// Process asset increase
	err := database.AddAssets(sessionID, outputs)
	if err != nil {
		slog.Error("ProcessExchangeResult", "error", "Failed to add assets", "err", err, "sessionID", sessionID)
		return err
	}

	slog.Info("ProcessExchangeResult", "sessionID", sessionID, "outputs", outputs, "action", "assets_added")
	return nil
}
