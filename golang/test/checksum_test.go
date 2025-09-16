package test

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// IsValidEVMAddress checks if a string is a valid EVM address
func IsValidEVMAddress(address string) bool {
	if address == "" {
		return false
	}

	// Remove 0x prefix if present
	cleanAddress := strings.TrimPrefix(address, "0x")

	// Check length (should be 40 hex characters)
	if len(cleanAddress) != 40 {
		return false
	}

	// Use go-ethereum's validation
	return common.IsHexAddress("0x" + cleanAddress)
}

// IsValidChecksumAddress checks if address has valid EIP-55 checksum
func IsValidChecksumAddress(address string) bool {
	if !IsValidEVMAddress(address) {
		return false
	}

	// Convert to common.Address and back to check checksum
	addr := common.HexToAddress(address)
	return addr.Hex() == address
}

func TestIsValidEVMAddress(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected bool
	}{
		// Valid addresses
		{
			name:     "Valid address with 0x prefix",
			address:  "0x742C4D61c2a093537BDBdea15Ec35b3C4094d0e3",
			expected: true,
		},
		{
			name:     "Valid address without 0x prefix",
			address:  "742C4D61c2a093537BDBdea15Ec35b3C4094d0e3",
			expected: true,
		},
		{
			name:     "Valid address all lowercase",
			address:  "0x742c4d61c2a093537bdbdea15ec35b3c4094d0e3",
			expected: true,
		},
		{
			name:     "Valid address all uppercase",
			address:  "0x742C4D61C2A093537BDBDEA15EC35B3C4094D0E3",
			expected: true,
		},
		{
			name:     "Zero address",
			address:  "0x0000000000000000000000000000000000000000",
			expected: true,
		},

		// Invalid addresses
		{
			name:     "Empty string",
			address:  "",
			expected: false,
		},
		{
			name:     "Too short",
			address:  "0x742C4D61c2a093537BDBdea15Ec35b3C4094d0e",
			expected: false,
		},
		{
			name:     "Too long",
			address:  "0x742C4D61c2a093537BDBdea15Ec35b3C4094d0e33a",
			expected: false,
		},
		{
			name:     "Invalid hex characters",
			address:  "0x742C4D61c2a093537BDBdea15Ec35b3C4094d0g3",
			expected: false,
		},
		{
			name:     "Missing characters",
			address:  "0x742C4D61c2a093537BDBdea15Ec35b3C4094d",
			expected: false,
		},
		{
			name:     "Only 0x prefix",
			address:  "0x",
			expected: false,
		},
		{
			name:     "Invalid characters in middle",
			address:  "0x742C4D61c2a093537BDBdXa15Ec35b3C4094d0e3",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEVMAddress(tt.address)
			assert.Equal(t, tt.expected, result, "Address: %s", tt.address)
		})
	}
}

func TestIsValidChecksumAddress(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected bool
	}{
		// Valid checksum addresses (EIP-55)
		{
			name:     "Valid checksum address",
			address:  "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
			expected: true,
		},
		{
			name:     "Valid checksum address 2",
			address:  "0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359",
			expected: true,
		},
		{
			name:     "Valid checksum address 3",
			address:  "0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB",
			expected: true,
		},
		{
			name:     "Valid checksum address 4",
			address:  "0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb",
			expected: true,
		},

		// Invalid checksum (wrong case)
		{
			name:     "Invalid checksum - wrong case",
			address:  "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", // all lowercase
			expected: false,
		},
		{
			name:     "Invalid checksum - wrong case 2",
			address:  "0x5AAEB6053F3E94C9B9A09F33669435E7EF1BEAED", // all uppercase
			expected: false,
		},
		{
			name:     "Invalid checksum - mixed wrong case",
			address:  "0x5aAeb6053f3e94c9b9a09f33669435e7ef1beaed",
			expected: false,
		},

		// All lowercase addresses (valid but no checksum)
		{
			name:     "All lowercase - no checksum validation needed",
			address:  "0x742c4d61c2a093537bdbdea15ec35b3c4094d0e3",
			expected: false, // Will fail checksum validation since it's not proper checksum format
		},

		// All uppercase addresses (valid but no checksum)
		{
			name:     "All uppercase - no checksum validation needed",
			address:  "0x742C4D61C2A093537BDBDEA15EC35B3C4094D0E3",
			expected: false, // Will fail checksum validation since it's not proper checksum format
		},

		// Invalid addresses
		{
			name:     "Invalid address format",
			address:  "invalid-address",
			expected: false,
		},
		{
			name:     "Empty string",
			address:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidChecksumAddress(tt.address)
			assert.Equal(t, tt.expected, result, "Address: %s", tt.address)
		})
	}
}

func TestCommonAddressOperations(t *testing.T) {
	t.Run("Convert string to common.Address and back", func(t *testing.T) {
		originalAddress := "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"

		// Convert to common.Address
		addr := common.HexToAddress(originalAddress)

		// Convert back to string with checksum
		convertedBack := addr.Hex()

		assert.Equal(t, originalAddress, convertedBack)
	})

	t.Run("Address comparison", func(t *testing.T) {
		addr1 := common.HexToAddress("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed")
		addr2 := common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed") // same address, different case

		// Addresses should be equal regardless of case
		assert.Equal(t, addr1, addr2)
	})

	t.Run("Zero address check", func(t *testing.T) {
		zeroAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
		nonZeroAddr := common.HexToAddress("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed")

		assert.Equal(t, common.Address{}, zeroAddr)
		assert.NotEqual(t, common.Address{}, nonZeroAddr)
	})
}
