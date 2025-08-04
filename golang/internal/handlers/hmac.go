package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// TODO: HMAC salt - In actual implementation, load from environment variables or configuration file
	HMACSalt = "my_secret_salt_value_!@#$%^&*"
)

// ValidateHMAC validates HMAC signature for request body
func ValidateHMAC(c *gin.Context) bool {
	// Get HMAC signature from header
	hmacSignature := c.GetHeader("X-HMAC-Signature")
	if hmacSignature == "" {
		return false
	}

	// Read request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return false
	}

	// Restore body for later use
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// Calculate HMAC
	hmacHash := hmac.New(sha256.New, []byte(HMACSalt))
	hmacHash.Write(bodyBytes)
	calculatedHash := hex.EncodeToString(hmacHash.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(calculatedHash), []byte(hmacSignature))
}

// HMACMiddleware middleware for HMAC validation
func HMACMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip HMAC validation for GET requests
		if c.Request.Method == "GET" {
			c.Next()
			return
		}

		// Validate HMAC for POST requests
		if !ValidateHMAC(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":   false,
				"errorCode": "INVALID_HMAC_SIGNATURE",
				"message":   "Invalid HMAC signature",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateHMAC generates HMAC signature for given data
func GenerateHMAC(data interface{}) (string, error) {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hmacHash := hmac.New(sha256.New, []byte(HMACSalt))
	hmacHash.Write(bodyBytes)
	return hex.EncodeToString(hmacHash.Sum(nil)), nil
}
