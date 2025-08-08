package models

import (
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

// ExchangeReq exchange response structure
type ExchangeReq struct {
	UUID    string           `json:"uuid"`
	TxHash  common.Hash      `json:"tx_hash"`
	Receipt ethTypes.Receipt `json:"receipt"`
	Intent  ExchangeIntent   `json:"intent"`
}

// PairAsset asset pair structure
type PairAsset struct {
	Type    string `json:"type"`
	AssetID string `json:"id"`
	Amount  uint   `json:"amount"`
}

// ExchangeIntent exchange intent structure
type ExchangeIntent struct {
	Type   string      `json:"type"`
	Method string      `json:"method"`
	From   []PairAsset `json:"from"`
	To     []PairAsset `json:"to"`
}
