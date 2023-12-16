package models

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// todo create better test data
	testPublicKey               = "key12345"
	testPublicKeyLastUpdateHash = "testHash1234"
)

// TestPublicKey will test public keys
func (ts *TestSuite) TestPublicKey() {
	ts.T().Run("success - no options, base model", func(t *testing.T) {
		key := NewPublicKey()
		require.NotNil(t, key)
		assert.NotNil(t, key.Logger())
		assert.Equal(t, uint64(0), key.GetID())
		assert.Equal(t, model.NamePublicKey.String(), key.Name())
		assert.Equal(t, model.TablePublicKeys, key.GetTableName())
	})

	ts.T().Run("success - create new public key", func(t *testing.T) {
		key := NewPublicKey(model.WithAllDependencies(ts.Dependencies), model.New())
		require.NotNil(t, key)
		assert.Equal(t, uint64(0), key.GetID())

		key.Key = testPublicKey
		key.LastUpdateHash = testPublicKeyLastUpdateHash
		key.Active = true

		err := key.Save(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(1), key.GetID())
		assert.Equal(t, testPublicKey, key.Key)
		assert.Equal(t, testPublicKeyLastUpdateHash, key.LastUpdateHash)
		assert.True(t, key.Active)
	})
}

// TestPublicKey_GetActiveKeys will test getting active keys
func (ts *TestSuite) TestPublicKey_GetActiveKeys() {

	// Create a new public key
	key := NewPublicKey(model.WithAllDependencies(ts.Dependencies), model.New())
	ts.Require().NotNil(key)
	ts.Require().Equal(uint64(0), key.GetID())

	key.Key = testPublicKey
	key.LastUpdateHash = testPublicKeyLastUpdateHash
	key.Active = true

	err := key.Save(context.Background())
	ts.Require().NoError(err)

	// Get the active keys
	var keys []*PublicKey
	keys, err = GetActivePublicKey(context.Background(), nil, model.WithAllDependencies(ts.Dependencies))
	ts.Require().NoError(err)
	ts.Require().NotNil(keys)
	ts.Require().Len(keys, 1)
	ts.Require().Equal(testPublicKey, keys[0].Key)
	ts.Require().Equal(testPublicKeyLastUpdateHash, keys[0].LastUpdateHash)
	ts.Require().True(keys[0].Active)
}

// TestPublicKey_GetActiveKeys_None will test getting active keys
func (ts *TestSuite) TestPublicKey_GetActiveKeys_None() {

	// Create a new public key
	key := NewPublicKey(model.WithAllDependencies(ts.Dependencies), model.New())
	ts.Require().NotNil(key)
	ts.Require().Equal(uint64(0), key.GetID())

	key.Key = testPublicKey
	key.LastUpdateHash = testPublicKeyLastUpdateHash
	key.Active = false

	err := key.Save(context.Background())
	ts.Require().NoError(err)

	// No active key found
	var keys []*PublicKey
	keys, err = GetActivePublicKey(context.Background(), nil, model.WithAllDependencies(ts.Dependencies))
	ts.Require().NoError(err)
	ts.Require().NotNil(keys)
	ts.Require().Empty(keys)
}
