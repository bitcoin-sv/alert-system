package models

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAlertMessageUnbanPeerRead tests the Read method of the AlertMessageUnbanPeer type
func TestAlertMessageUnbanPeerRead(t *testing.T) {

	t.Run("valid unban peer alert just IP", func(t *testing.T) {
		// Create a sample alert message payload
		peer := []byte("127.0.0.1:8333")
		reason := []byte("test")

		// Encode the lengths of peer and reason using binary.Write
		payload := make([]byte, 0)
		payload = append(payload, encodeVarInt(uint64(len(peer)))...)
		payload = append(payload, peer...)
		payload = append(payload, encodeVarInt(uint64(len(reason)))...)
		payload = append(payload, reason...)

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString(hex.EncodeToString(payload))
		require.NoError(t, err)

		// Create an AlertMessageUnbanPeer instance
		alert := &AlertMessageUnbanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.NoError(t, err)
		assert.Equal(t, uint64(len(peer)), alert.PeerLength)
		assert.Equal(t, peer, alert.Peer)
		assert.Equal(t, uint64(len(reason)), alert.ReasonLength)
		assert.Equal(t, reason, alert.Reason)
	})

	t.Run("valid unban peer alert IP and port", func(t *testing.T) {
		// Create a sample alert message payload for unban peer with subnet
		peer := []byte("127.0.0.1/24")
		reason := []byte("test")

		// Encode the lengths of peer and reason using binary.Write
		/*payload := make([]byte, 0)
		payload = append(payload, encodeVarInt(uint64(len(peer)))...)
		payload = append(payload, peer...)
		payload = append(payload, encodeVarInt(uint64(len(reason)))...)
		payload = append(payload, reason...)*/

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString("0c3132372e302e302e312f32340474657374")
		require.NoError(t, err)

		// Create an AlertMessageUnbanPeer instance
		alert := &AlertMessageUnbanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.NoError(t, err)
		assert.Equal(t, uint64(len(peer)), alert.PeerLength)
		assert.Equal(t, peer, alert.Peer)
		assert.Equal(t, uint64(len(reason)), alert.ReasonLength)
		assert.Equal(t, reason, alert.Reason)
	})

	t.Run("bad peer length, EOF", func(t *testing.T) {
		// Create a sample alert message payload with mismatched PeerLength
		// peer := []byte("127.0.0.1/24")
		// reason := []byte("test")

		// Encode the lengths of peer and reason using binary.Write
		/*payload := make([]byte, 0)
		payload = append(payload, encodeVarInt(uint64(len(peer)))...)
		payload = append(payload, peer...)
		payload = append(payload, encodeVarInt(uint64(len(reason)))...)
		payload = append(payload, reason...)*/

		// Add extra bytes to the payload to cause an EOF error
		// payload = append(payload, []byte{0x1c}...)

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString("1c3132372e302e302e312f32340474657374")
		require.NoError(t, err)

		// Create an AlertMessageUnbanPeer instance
		alert := &AlertMessageUnbanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.Error(t, err) // Expects an error due to mismatched PeerLength and EOF
	})

	t.Run("bad peer length, bad data", func(t *testing.T) {
		// Create a sample alert message payload with mismatched PeerLength and bad data
		peer := []byte("127.0.0.1/24")
		reason := []byte("test")

		// Encode the lengths of peer and reason using binary.Write
		payload := make([]byte, 0)
		payload = append(payload, encodeVarInt(uint64(len(peer)))...)
		payload = append(payload, peer...)
		payload = append(payload, encodeVarInt(uint64(len(reason)))...)
		payload = append(payload, reason...)

		// Replace the PeerLength bytes with incorrect data
		payload[0] = 0x02

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString("023132372e302e302e312f32340474657374")
		require.NoError(t, err)

		// Create an AlertMessageUnbanPeer instance
		alert := &AlertMessageUnbanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.Error(t, err) // Expects an error due to mismatched PeerLength and bad data
	})

	t.Run("nonsensical data", func(t *testing.T) {
		// Create a sample alert message payload with nonsensical data
		alertData := "02311123491827331928437cdf283721"

		// Convert the args.alert string to bytes
		alertBytes, err := hex.DecodeString(alertData)
		require.NoError(t, err)

		// Create an AlertMessageUnbanPeer instance
		alert := &AlertMessageUnbanPeer{}

		// Call the Read method to parse the payload
		err = alert.Read(alertBytes)

		// Use the testify/assert package for assertions
		require.Error(t, err) // Expects an error due to nonsensical data
	})
}
