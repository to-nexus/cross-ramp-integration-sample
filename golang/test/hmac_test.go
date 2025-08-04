package test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	// TODO: HMAC salt - In actual implementation, load from environment variables or configuration file
	salt = "my_secret_salt_value_!@#$%^&*" // hmac key
)

type Body struct {
	UserID    int    `json:"userId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt int    `json:"createdAt"`
}

var (
	body = Body{
		UserID:    1234,
		Username:  "홍길동",
		Email:     "user@example.com",
		Role:      "admin",
		CreatedAt: 1234567890,
	}
)

func TestSha256(t *testing.T) {
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)
	t.Log(string(bodyBytes))

	hmac := hmac.New(sha256.New, []byte(salt))
	hmac.Write(bodyBytes)
	hashBytes := hmac.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	t.Log("hashString", hashString) // expected X-HMAC-Signature: f96cf60394f6b8ad3c6de2d5b2b1d1a540f9529082a8eb9cee405bfbdd9f37a1
}
