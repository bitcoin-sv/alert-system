package models

import (
	"encoding/hex"
	"github.com/bsv-blockchain/go-sdk/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAlertMessageInvalidateBlock_Read will test the method Read()
func TestAlertMessageInvalidateBlock_Read(t *testing.T) {
	tests := []struct {
		name        string
		blockHash   string // hex string
		reason      []byte
		expectError bool
	}{
		{
			name:      "Successful Read of block hash to invalidate",
			blockHash: "00000000000000000ab414417d6197f620a3917dc25d6fac7191de37739c45d6",
			reason: []byte{
				0x05, 'h', 'e', 'l', 'l', 'o',
			},
			expectError: false,
		},
		{
			name:      "Unsuccessful Read of block hash to invalidate; wrong reason length",
			blockHash: "00000000000000000ab414417d6197f620a3917dc25d6fac7191de37739c45d6",
			reason: []byte{
				0x04, 'h', 'e', 'l', 'l', 'o',
			},
			expectError: true,
		},
		{
			name:      "Unsuccessful Read of block hash to invalidate; block hash too short",
			blockHash: "00000000000000000ab414417d6197f620a3917dc25d6fac7191de37739c45",
			reason: []byte{
				0x05, 'h', 'e', 'l', 'l', 'o',
			},
			expectError: true,
		},
		{
			name:      "Unsuccessful Read of block hash to invalidate; block hash too long",
			blockHash: "00000000000000000ab414417d6197f620a3917dc25d6fac7191de37739c45d655",
			reason: []byte{
				0x05, 'h', 'e', 'l', 'l', 'o',
			},
			expectError: true,
		},
		{
			name:        "Unsuccessful Read of block hash to invalidate; no reason",
			blockHash:   "00000000000000000ab414417d6197f620a3917dc25d6fac7191de37739c45d655",
			reason:      []byte{},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlertMessageInvalidateBlock{}
			alertBytes, err := hex.DecodeString(tt.blockHash)
			require.NoError(t, err)
			alertBytes = util.ReverseBytes(alertBytes)
			alertBytes = append(alertBytes, tt.reason...)
			err = a.Read(alertBytes)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.blockHash, a.BlockHash.String())
				assert.Equal(t, uint64(len(tt.reason)-1), a.ReasonLength)
				assert.Equal(t, tt.reason[1:], a.Reason)
			}
		})
	}
}
