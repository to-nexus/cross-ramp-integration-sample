package test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestCryptoSign(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Fatalf("Failed to convert hex to ECDSA: %v", err)
	}

	digest := crypto.Keccak256([]byte("test1"))
	signature, err := crypto.Sign(digest, privateKey)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	t.Log("digest: ", hexutil.Encode(digest))
	t.Log("signature: ", hexutil.Encode(signature))
	t.Log("signature size: ", len(signature))
}

func TestCryptoSign2(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Fatalf("Failed to convert hex to ECDSA: %v", err)
	}

	digest := common.Hex2Bytes("d91c81e564e4f69229a9224943fa9a79ff21b60fcef5096bfb79e1ce28591a85")
	signature, err := crypto.Sign(digest, privateKey)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	t.Log(hexutil.Encode(signature))
}
