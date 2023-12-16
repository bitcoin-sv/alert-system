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

func TestAlertMessageInformational_MessageString(t *testing.T) {
	type fields struct {
		AlertMessage  AlertMessage
		MessageLength uint64
		Message       []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{{
		name: "test valid message",
		fields: fields{
			AlertMessage: AlertMessage{
				Raw: "010000001b00000067b5bd6500000000010000000774657374696e67202214d4892217b450eedfb33dd901951e80557ea10d19a59f8a566f733b1ab7107b77d388a9f2fac6602b7258cbcb0ac11c9a6dd0b5687cb9508bcfa5dbd6ce901f4672d99e36978856f2d2794c4c48d353a0b45357d08991147f9e8803a0b90a5f01e85739f36eab32765fe2190b1625e3f5d6c41319da3da803b60be472bf2c011f3784e3d3504c93be28e32e9108aead94cb515bb4813303e6a14735bcca87e451487b222198a9ba3ea0c984e3fbd95e35ba1607c5c74224af6083185a17ea7ff9",
			},
			Message: []byte("testing"),
		},
		want: "Informational: testing",
	},

	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlertMessageInformational{
				AlertMessage: tt.fields.AlertMessage,
				//MessageLength: tt.fields.MessageLength,
				Message: tt.fields.Message,
			}
			assert.Equalf(t, tt.want, a.MessageString(), "MessageString()")
		})
	}
}
