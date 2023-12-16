package models

import (
	"encoding/hex"
	"testing"

	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMessageBanPeer_Read tests the Read method of the MessageBanPeer struct
func (ts *TestSuite) TestMessageBanPeer_Read() {
	type fields struct {
		PeerLength   uint64
		Peer         []byte
		ReasonLength uint64
		Reason       []byte
	}
	type args struct {
		alert string // hex string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid ban peer alert just IP",
			fields: fields{
				PeerLength:   uint64(9),
				Peer:         []byte("127.0.0.1"),
				ReasonLength: uint64(4),
				Reason:       []byte("test"),
			},
			args: args{
				alert: "093132372e302e302e310474657374",
			},
		},
		{
			name: "valid ban peer alert with subnet",
			fields: fields{
				PeerLength:   uint64(12),
				Peer:         []byte("127.0.0.1/24"),
				ReasonLength: uint64(4),
				Reason:       []byte("test"),
			},
			args: args{
				alert: "0c3132372e302e302e312f32340474657374",
			},
		},
		{
			name: "bad peer length, EOF",
			fields: fields{
				PeerLength:   uint64(12),
				Peer:         []byte("127.0.0.1/24"),
				ReasonLength: uint64(4),
				Reason:       []byte("test"),
			},
			args: args{
				alert: "1c3132372e302e302e312f32340474657374",
			},
			wantErr: true,
		},
		{
			name: "bad peer length, bad data",
			fields: fields{
				PeerLength:   uint64(12),
				Peer:         []byte("127.0.0.1/24"),
				ReasonLength: uint64(4),
				Reason:       []byte("test"),
			},
			args: args{
				alert: "023132372e302e302e312f32340474657374",
			},
			wantErr: true,
		},
		{
			name: "nonsensical data",
			fields: fields{
				PeerLength:   uint64(12),
				Peer:         []byte("127.0.0.1/24"),
				ReasonLength: uint64(4),
				Reason:       []byte("test"),
			},
			args: args{
				alert: "02311123491827331928437cdf283721",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			alertBytes, err := hex.DecodeString(tt.args.alert)
			if err != nil {
				t.Errorf("invalid hex string: %v", err.Error())
				return
			}
			alert := NewAlertMessage(model.WithAllDependencies(ts.Dependencies), model.New())
			alert.SetAlertType(AlertTypeBanPeer)
			a := alert.ProcessAlertMessage()
			if err = a.Read(alertBytes); (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
			/*if !bytes.Equal(a.Peer, tt.fields.Peer) && !tt.wantErr {
				t.Errorf("peer [%s] does not match wanted [%s]", a.Peer, tt.fields.Peer)
			}
			if !bytes.Equal(a.Reason, tt.fields.Reason) && !tt.wantErr {
				t.Errorf("reason [%s] does not match wanted [%s]", a.Reason, tt.fields.Reason)
			}*/
		})
	}
}

// TestAlertMessageBanPeerRead tests the Read method of the AlertMessageBanPeer struct
func TestAlertMessageBanPeerRead(t *testing.T) {

	t.Run("valid ban peer alert just IP", func(t *testing.T) {
		// Create a sample alert message payload
		peerLength := uint64(9)
		peer := []byte("127.0.0.1")
		reasonLength := uint64(4)
		reason := []byte("test")

		// Encode the lengths of peer and reason using binary.Write
		//payload := make([]byte, 0)
		//payload = append(payload, encodeVarInt(peerLength)...)
		//payload = append(payload, peer...)
		//payload = append(payload, encodeVarInt(reasonLength)...)
		//payload = append(payload, reason...)

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString("093132372e302e302e310474657374")
		require.NoError(t, err)

		// Create an AlertMessageBanPeer instance
		alert := &AlertMessageBanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.NoError(t, err)
		assert.Equal(t, peerLength, alert.PeerLength)
		assert.Equal(t, peer, alert.Peer)
		assert.Equal(t, reasonLength, alert.ReasonLength)
		assert.Equal(t, reason, alert.Reason)
	})

	t.Run("valid ban peer alert with subnet", func(t *testing.T) {
		// Create a sample alert message payload with subnet
		peerLength := uint64(12)
		peer := []byte("127.0.0.1/24")
		reasonLength := uint64(4)
		reason := []byte("test")

		// Encode the lengths of peer and reason using binary.Write
		/*payload := make([]byte, 0)
		payload = append(payload, encodeVarInt(peerLength)...)
		payload = append(payload, peer...)
		payload = append(payload, encodeVarInt(reasonLength)...)
		payload = append(payload, reason...)*/

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString("0c3132372e302e302e312f32340474657374")
		require.NoError(t, err)

		// Create an AlertMessageBanPeer instance
		alert := &AlertMessageBanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.NoError(t, err)
		assert.Equal(t, peerLength, alert.PeerLength)
		assert.Equal(t, peer, alert.Peer)
		assert.Equal(t, reasonLength, alert.ReasonLength)
		assert.Equal(t, reason, alert.Reason)
	})

	t.Run("bad peer length, EOF", func(t *testing.T) {
		// Create a sample alert message payload with a mismatched PeerLength
		// peerLength := uint64(12) // Mismatched with actual Peer length
		// peer := []byte("127.0.0.1/24")
		// reasonLength := uint64(4)
		// reason := []byte("test")

		// Encode the lengths of peer and reason using binary.Write
		/*payload := make([]byte, 0)
		payload = append(payload, encodeVarInt(peerLength)...)
		payload = append(payload, peer...)
		payload = append(payload, encodeVarInt(reasonLength)...)
		payload = append(payload, reason...)*/

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString("1c3132372e302e302e312f32340474657374")
		require.NoError(t, err)

		// Create an AlertMessageBanPeer instance
		alert := &AlertMessageBanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.Error(t, err) // Expects an error due to mismatched PeerLength
	})

	t.Run("bad peer length, bad data", func(t *testing.T) {
		// Create a sample alert message payload with a mismatched PeerLength and bad peer data
		// peerLength := uint64(12) // Mismatched with actual Peer length
		// peer := []byte("127.0.0.1/24")
		// reasonLength := uint64(4)
		// reason := []byte("test")

		// Encode the lengths of peer and reason using binary.Write
		/*payload := make([]byte, 0)
		payload = append(payload, encodeVarInt(peerLength)...)
		// Incorrect peer data length (2 bytes) and incorrect data
		payload = append(payload, []byte{0x02, 0x31}...)
		payload = append(payload, encodeVarInt(reasonLength)...)
		payload = append(payload, reason...)*/

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString("023132372e302e302e312f32340474657374")
		require.NoError(t, err)

		// Create an AlertMessageBanPeer instance
		alert := &AlertMessageBanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.Error(t, err) // Expects an error due to mismatched PeerLength and bad peer data
	})

	t.Run("nonsensical data", func(t *testing.T) {
		// Create a sample alert message payload with nonsensical data
		// peerLength := uint64(12)
		// peer := []byte("127.0.0.1/24")
		// reasonLength := uint64(4)
		// reason := []byte("test")

		// Encode the lengths of peer and reason using binary.Write
		// payload := make([]byte, 0)
		// payload = append(payload, encodeVarInt(peerLength)...)
		// payload = append(payload, peer...)
		// payload = append(payload, encodeVarInt(reasonLength)...)
		// payload = append(payload, reason...)

		// Add extra nonsensical data to the payload
		// payload = append(payload, []byte{0x23, 0x11, 0x12, 0x34, 0x91, 0x82, 0x73, 0x31, 0x92, 0x84, 0x37, 0xcd, 0xf2, 0x83, 0x72, 0x1}...)

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString("02311123491827331928437cdf283721")
		require.NoError(t, err)

		// Create an AlertMessageBanPeer instance
		alert := &AlertMessageBanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.Error(t, err) // Expects an error due to nonsensical data
	})
}

// encodeVarInt encodes a uint64 into a variable length byte slice
func encodeVarInt(value uint64) []byte {
	var buf []byte
	for {
		b := byte(value & 0x7f)
		value >>= 7
		if value != 0 {
			b |= 0x80
		}
		buf = append(buf, b)
		if value == 0 {
			break
		}
	}
	return buf
}
