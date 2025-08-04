package models

import (
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

// ExchangeResultRequest exchange result request structure
type ExchangeResultRequest struct {
	UUID    string           `json:"uuid"`
	TxHash  common.Hash      `json:"tx_hash"`
	Receipt ethTypes.Receipt `json:"receipt"`
	Intent  ExchangeIntent   `json:"intent"`
}

// PairToken token pair structure
type PairToken struct {
	TokenID string `json:"id"`
	Amount  uint   `json:"amount"`
}

// PairAsset asset pair structure
type PairAsset struct {
	AssetID string `json:"id"`
	Amount  uint   `json:"amount"`
}

// ExchangeIntent exchange intent structure
type ExchangeIntent struct {
	ProjectID string      `json:"project_id"`
	PairID    uint        `json:"pair_id"`
	Token     PairToken   `json:"token"`
	Materials []PairAsset `json:"materials,omitempty"`
	Outputs   []PairAsset `json:"outputs,omitempty"`
}
