package database

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"sample-game-backend/internal/models"

	"github.com/hashicorp/go-memdb"
)

// DB global variables
var (
	db     *memdb.MemDB
	dbOnce sync.Once
	dbInit bool
)

// InitDB initialize database (singleton pattern)
func InitDB() error {
	var initErr error
	dbOnce.Do(func() {
		// Schema definition
		schema := &memdb.DBSchema{
			Tables: map[string]*memdb.TableSchema{
				"session_assets": {
					Name: "session_assets",
					Indexes: map[string]*memdb.IndexSchema{
						"id": {
							Name:    "id",
							Unique:  true,
							Indexer: &memdb.StringFieldIndex{Field: "SessionID"},
						},
					},
				},
				"uuid_mapping": {
					Name: "uuid_mapping",
					Indexes: map[string]*memdb.IndexSchema{
						"id": {
							Name:    "id",
							Unique:  true,
							Indexer: &memdb.StringFieldIndex{Field: "UUID"},
						},
					},
				},
			},
		}

		db, initErr = memdb.NewMemDB(schema)
		if initErr == nil {
			dbInit = true
			slog.Info("InitDB", "status", "success", "message", "Database initialized successfully")
		} else {
			slog.Error("InitDB", "error", "Failed to initialize database", "err", initErr)
		}
	})

	return initErr
}

// GetDB return database instance (check initialization)
func GetDB() (*memdb.MemDB, error) {
	if !dbInit {
		return nil, fmt.Errorf("database not initialized")
	}
	return db, nil
}

// CloseDB close database connection
func CloseDB() error {
	// memdb is in-memory, so no additional cleanup is needed
	return nil
}

// generateRandomAssets generate random assets
func generateRandomAssets() map[string]string {
	assets := make(map[string]string)
	baseAmount := 100000000
	// Generate asset_money randomly (1000 ~ 5000)
	moneyAmount := rand.Intn(baseAmount) + 1000
	assets["asset_money"] = strconv.Itoa(moneyAmount)

	// Generate asset_gold randomly (500 ~ 3000)
	goldAmount := rand.Intn(baseAmount) + 500
	assets["asset_gold"] = strconv.Itoa(goldAmount)

	// Generate item_gem randomly (500 ~ 3000)
	gemAmount := rand.Intn(baseAmount) + 500
	assets["item_gem"] = strconv.Itoa(gemAmount)

	// Generate item_banana randomly (500 ~ 3000)
	bananaAmount := rand.Intn(baseAmount) + 500
	assets["item_banana"] = strconv.Itoa(bananaAmount)

	// Generate asset_silver randomly (500 ~ 3000)
	silverAmount := rand.Intn(baseAmount) + 500
	assets["asset_silver"] = strconv.Itoa(silverAmount)

	// Generate item_apple randomly (500 ~ 3000)
	appleAmount := rand.Intn(baseAmount) + 500
	assets["item_apple"] = strconv.Itoa(appleAmount)

	// Generate item_fish randomly (500 ~ 3000)
	fishAmount := rand.Intn(baseAmount) + 500
	assets["item_fish"] = strconv.Itoa(fishAmount)

	// Generate item_branch randomly (500 ~ 3000)
	branchAmount := rand.Intn(baseAmount) + 500
	assets["item_branch"] = strconv.Itoa(branchAmount)

	// Generate item_horn randomly (500 ~ 3000)
	hornAmount := rand.Intn(baseAmount) + 500
	assets["item_horn"] = strconv.Itoa(hornAmount)

	// Generate item_maple randomly (500 ~ 3000)
	mapleAmount := rand.Intn(baseAmount) + 500
	assets["item_maple"] = strconv.Itoa(mapleAmount)

	return assets
}

// GetOrCreateSessionAssets get or create session-specific asset information
func GetOrCreateSessionAssets(sessionID string) (*models.SessionAssets, error) {
	database, err := GetDB()
	if err != nil {
		return nil, err
	}

	// Start read transaction
	txn := database.Txn(false)
	defer txn.Abort()

	// Query existing data
	raw, err := txn.First("session_assets", "id", sessionID)
	if err != nil {
		return nil, err
	}

	if raw != nil {
		// Return existing data if found
		sessionAssets := raw.(*models.SessionAssets)
		return sessionAssets, nil
	}

	// Create new session assets
	sessionAssets := &models.SessionAssets{
		SessionID: sessionID,
		Assets:    generateRandomAssets(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Start write transaction
	txn = database.Txn(true)
	err = txn.Insert("session_assets", sessionAssets)
	if err != nil {
		txn.Abort()
		return nil, err
	}

	txn.Commit()
	return sessionAssets, nil
}

// CheckAndDeductAssets validate and deduct asset balance
func CheckAndDeductAssets(sessionID string, fromAssets []struct {
	Type   string `json:"type" binding:"required"`
	ID     string `json:"id" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}) error {
	// Get session asset information
	sessionAssets, err := GetOrCreateSessionAssets(sessionID)
	if err != nil {
		return err
	}

	// Validate and deduct balance for each asset
	for _, asset := range fromAssets {
		currentBalance, exists := sessionAssets.Assets[asset.ID]
		if !exists {
			return fmt.Errorf("asset %s not found in session", asset.ID)
		}

		// Convert string to integer
		currentAmount, err := strconv.Atoi(currentBalance)
		if err != nil {
			return fmt.Errorf("invalid balance format for asset %s", asset.ID)
		}

		// Validate balance
		if currentAmount < asset.Amount {
			return fmt.Errorf("insufficient balance for asset %s: required %d, available %d", asset.ID, asset.Amount, currentAmount)
		}

		// Deduct
		newBalance := currentAmount - asset.Amount
		sessionAssets.Assets[asset.ID] = strconv.Itoa(newBalance)

	}

	// Set update time
	sessionAssets.UpdatedAt = time.Now().Format(time.RFC3339)

	// Save to DB
	txn := db.Txn(true)
	err = txn.Insert("session_assets", sessionAssets)
	if err != nil {
		txn.Abort()
		return err
	}

	txn.Commit()
	return nil
}

// AddAssets increase assets
func AddAssets(sessionID string, assets []models.PairAsset) error {
	// Get session asset information
	sessionAssets, err := GetOrCreateSessionAssets(sessionID)
	if err != nil {
		return err
	}

	// Increase balance for each asset
	for _, asset := range assets {
		currentBalance, exists := sessionAssets.Assets[asset.AssetID]
		if !exists {
			// Create new asset if it doesn't exist
			sessionAssets.Assets[asset.AssetID] = strconv.FormatUint(uint64(asset.Amount), 10)
		} else {
			// Add to existing balance
			currentAmount, err := strconv.ParseUint(currentBalance, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid balance format for asset %s", asset.AssetID)
			}

			newBalance := currentAmount + uint64(asset.Amount)
			sessionAssets.Assets[asset.AssetID] = strconv.FormatUint(newBalance, 10)
		}
	}

	// Set update time
	sessionAssets.UpdatedAt = time.Now().Format(time.RFC3339)

	// Save to DB
	txn := db.Txn(true)
	err = txn.Insert("session_assets", sessionAssets)
	if err != nil {
		txn.Abort()
		return err
	}

	txn.Commit()
	return nil
}

// UUIDMapping UUID 매핑 구조체
type UUIDMapping struct {
	UUID      string `json:"uuid"`
	SessionID string `json:"session_id"`
}

// StoreUUIDMapping UUID와 SessionID 매핑 저장
func StoreUUIDMapping(uuid, sessionID string) error {
	mapping := &UUIDMapping{
		UUID:      uuid,
		SessionID: sessionID,
	}

	database, err := GetDB()
	if err != nil {
		slog.Error("StoreUUIDMapping", "error", "Failed to get database", "err", err)
		return err
	}

	txn := database.Txn(true)
	err = txn.Insert("uuid_mapping", mapping)
	if err != nil {
		txn.Abort()
		return err
	}

	txn.Commit()
	slog.Info("StoreUUIDMapping", "uuid", uuid, "sessionID", sessionID, "action", "committed")
	return nil
}

// GetSessionIDByUUID UUID로 SessionID 조회
func GetSessionIDByUUID(uuid string) (string, error) {
	database, err := GetDB()
	if err != nil {
		return "", err
	}

	txn := database.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("uuid_mapping", "id", uuid)
	if err != nil {
		return "", err
	}

	if raw == nil {
		slog.Warn("GetSessionIDByUUID", "warning", "UUID not found", "uuid", uuid)
		return "", fmt.Errorf("uuid mapping not found: %s", uuid)
	}

	mapping := raw.(*UUIDMapping)
	sessionID := mapping.SessionID
	slog.Info("GetSessionIDByUUID", "uuid", uuid, "sessionID", sessionID)
	return sessionID, nil
}
