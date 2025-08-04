package models

// V1Data v1 guide data structure
type V1Data struct {
	PlayerID      string  `json:"player_id"`
	Name          string  `json:"name"`
	WalletAddress string  `json:"wallet_address"`
	Server        string  `json:"server"`
	Assets        []Asset `json:"assets"`
}

// Response API response structure
type Response struct {
	Success   bool    `json:"success"`
	ErrorCode *string `json:"errorCode,omitempty"`
	Data      struct {
		V1    V1Data `json:"v1"`
		Guide any    `json:"guide"`
	} `json:"data"`
}

// ValidateRequest order validation request structure
type ValidateRequest struct {
	UUID        string `json:"uuid" binding:"required"`
	UserSig     string `json:"user_sig" binding:"required"`
	UserAddress string `json:"user_address" binding:"required"`
	ProjectID   string `json:"project_id" binding:"required"`
	Digest      string `json:"digest" binding:"required"`
	Intent      struct {
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
	} `json:"intent" binding:"required"`
}

// ValidateResponse order validation response structure
type ValidateResponse struct {
	Success   bool    `json:"success"`
	ErrorCode *string `json:"errorCode,omitempty"`
	Data      struct {
		UserSig      string `json:"userSig"`
		ValidatorSig string `json:"validatorSig"`
	} `json:"data"`
}
