package utils

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSignWithGenesis_Success tests the function SignWithGenesis
func TestSignWithGenesis_Success(t *testing.T) {
	// Prepare the input data
	data := []byte("test data")

	// Call the function
	signatures, err := SignWithGenesis(data)

	// Validate the results
	require.NoError(t, err)
	assert.Len(t, signatures, 3) // Since we expect 3 signatures

	// Further validation can be done depending on the expected output format
	assert.Equal(t, "1f4132a89ccf82d3df3a0ee06105eeee173c374930f84fe0c9f6587790b2cbc6ef65a33bbe4e4b9fe57c65ba816db2d62cff4c9b4e66531cfc549d2be746f495bb", hex.EncodeToString(signatures[0]))
	assert.Equal(t, "1f71934b00f38fc068ce606b9ea6c9a3a73b7f68ed2ec7cfea3ab93cf3b1e117c0605b74b026a9e575e1b63d55b54bd50f493d52cdd37fbaa7e0aba26b96dc0376", hex.EncodeToString(signatures[1]))
	assert.Equal(t, "1f17850af80b856545e77b060057e4c31a45bdc5fa89a6f91a6bb7c8df683b5875269d18ec3dc61c7198ff07d53df241fa61210f3825a871bd485e32329e6fcff2", hex.EncodeToString(signatures[2]))
}

// TestSignWithKeys_Success tests the function SignWithKeys
func TestSignWithKeys_Success(t *testing.T) {
	// Prepare the input data
	data := []byte("test data")
	keys := []string{Key1, Key2, Key3}

	// Call the function
	signatures, err := SignWithKeys(data, keys)

	// Validate the results
	require.NoError(t, err)
	assert.Len(t, signatures, len(keys)) // Length of signatures should match the number of keys

	// Further validation can be done depending on the expected output format
	assert.Equal(t, "1f4132a89ccf82d3df3a0ee06105eeee173c374930f84fe0c9f6587790b2cbc6ef65a33bbe4e4b9fe57c65ba816db2d62cff4c9b4e66531cfc549d2be746f495bb", hex.EncodeToString(signatures[0]))
	assert.Equal(t, "1f71934b00f38fc068ce606b9ea6c9a3a73b7f68ed2ec7cfea3ab93cf3b1e117c0605b74b026a9e575e1b63d55b54bd50f493d52cdd37fbaa7e0aba26b96dc0376", hex.EncodeToString(signatures[1]))
	assert.Equal(t, "1f17850af80b856545e77b060057e4c31a45bdc5fa89a6f91a6bb7c8df683b5875269d18ec3dc61c7198ff07d53df241fa61210f3825a871bd485e32329e6fcff2", hex.EncodeToString(signatures[2]))
}

// BenchmarkSignWithKeys benchmarks the function SignWithKeys
func BenchmarkSignWithKeys(b *testing.B) {
	data := []byte("test data")
	keys := []string{Key1, Key2, Key3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SignWithKeys(data, keys)
	}
}

// BenchmarkSignWithGenesis benchmarks the function SignWithGenesis
func BenchmarkSignWithGenesis(b *testing.B) {
	data := []byte("test data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SignWithGenesis(data)
	}
}
