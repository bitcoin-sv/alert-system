package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAlertMessageInformational_Read will test the method Read()
func TestAlertMessageInformational_Read(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectError bool
		expectedMsg AlertMessageInformational
	}{
		{
			name:        "Successful Read",
			input:       []byte{0x05, 'h', 'e', 'l', 'l', 'o'}, // 0x05 is the length prefix for "hello"
			expectError: false,
			expectedMsg: AlertMessageInformational{
				MessageLength: 5,
				Message:       []byte("hello"),
			},
		},
		{
			name:        "Error - Length Longer Than Buffer",
			input:       []byte{0x06, 'w', 'o', 'r', 'l', 'd'}, // 0x06 indicates a length of 6, but "world" is only 5 bytes
			expectError: true,
		},
		{
			name:        "Error - Invalid VarInt",
			input:       []byte{0xFF}, // 0xFF is not a valid VarInt for length
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlertMessageInformational{}
			err := a.Read(tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedMsg.MessageLength, a.MessageLength)
				assert.Equal(t, tt.expectedMsg.Message, a.Message)
			}
		})
	}
}
