package models

// Asset asset information structure
type Asset struct {
	ID      string `json:"id"`
	Balance string `json:"balance"`
}

// SessionAssets session-specific asset information structure
type SessionAssets struct {
	SessionID string            `json:"session_id"`
	Assets    map[string]string `json:"assets"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}
