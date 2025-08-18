package test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestCryptoSign(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Fatalf("Failed to convert hex to ECDSA: %v", err)
	}

	digest := crypto.Keccak256([]byte("test"))
	signature, err := crypto.Sign(digest, privateKey)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	t.Log(hexutil.Encode(signature))
}
