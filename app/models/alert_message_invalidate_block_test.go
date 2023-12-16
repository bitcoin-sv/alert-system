package models

/*func TestAlertMessageInvalidateBlockRead(t *testing.T) {
	// Create a sample alert message payload
	blockHashBytes, _ := hex.DecodeString("00000000000000000ab414417d6197f620a3917dc25d6fac7191de37739c45d6")
	reasonLength := uint64(8)
	reason := []byte("invalid")

	// Encode the lengths of reason using binary.Write
	payload := make([]byte, 0)
	payload = append(payload, blockHashBytes...)
	payload = append(payload, encodeVarInt(reasonLength)...)
	payload = append(payload, reason...)

	// Create an AlertMessageInvalidateBlock instance
	alert := &AlertMessageInvalidateBlock{}

	// Call the Read method to parse the payload
	err := alert.Read(payload)

	// Use the testify/assert package for assertions
	assert.NoError(t, err)
	// 	assert.Equal(t, blockHashBytes, alert.BlockHash.CloneBytes())
	assert.Equal(t, reasonLength, alert.ReasonLength)
	assert.Equal(t, reason, alert.Reason)
}*/
