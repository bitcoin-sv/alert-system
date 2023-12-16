package models

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// todo create better test data
	testAlertHash = "hash12345"
	testAlertRaw  = "raw_message"
)

// TestAlertMessage will test alert messages
func (ts *TestSuite) TestAlertMessage() {
	ts.T().Run("success - no options, base model", func(t *testing.T) {
		message := NewAlertMessage()
		require.NotNil(t, message)
		assert.NotNil(t, message.Logger())
		assert.Equal(t, uint64(0), message.GetID())
		assert.Equal(t, model.NameAlertMessage.String(), message.Name())
		assert.Equal(t, model.TableAlertMessages, message.GetTableName())
	})

	ts.T().Run("success - create new alert message", func(t *testing.T) {
		message := NewAlertMessage(model.WithAllDependencies(ts.Dependencies), model.New())
		require.NotNil(t, message)
		assert.Equal(t, uint64(0), message.GetID())

		message.Hash = testAlertHash
		message.Raw = testAlertRaw
		message.SequenceNumber = 1

		err := message.Save(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(1), message.GetID())
		assert.Equal(t, testAlertHash, message.Hash)
		assert.Equal(t, testAlertRaw, message.Raw)
		assert.Equal(t, uint32(1), message.SequenceNumber)
	})
}

// TestAlertMessage_GetAlertBySequenceNumber will test getting an alert by sequence number
func (ts *TestSuite) TestAlertMessage_GetAlertBySequenceNumber() {

	// Create a new alert message
	message := NewAlertMessage(model.WithAllDependencies(ts.Dependencies), model.New())
	ts.Require().NotNil(message)
	ts.Require().Equal(uint64(0), message.GetID())

	message.Hash = testAlertHash
	message.Raw = testAlertRaw
	message.SequenceNumber = 1

	err := message.Save(context.Background())
	ts.Require().NoError(err)

	// Get the alert message
	message, err = GetAlertMessageBySequenceNumber(context.Background(), 1, model.WithAllDependencies(ts.Dependencies))
	ts.Require().NoError(err)

	ts.Require().Equal(uint64(1), message.GetID())
	ts.Require().Equal(testAlertHash, message.Hash)
	ts.Require().Equal(testAlertRaw, message.Raw)
	ts.Require().Equal(uint32(1), message.SequenceNumber)
}

// TestAlertMessage_GetLatestAlert will test getting the latest alert
func (ts *TestSuite) TestAlertMessage_GetLatestAlert() {

	// Create the first alert message
	message := NewAlertMessage(model.WithAllDependencies(ts.Dependencies), model.New())
	ts.Require().NotNil(message)
	ts.Require().Equal(uint64(0), message.GetID())

	message.Hash = testAlertHash
	message.Raw = testAlertRaw
	message.SequenceNumber = 1

	err := message.Save(context.Background())
	ts.Require().NoError(err)

	// Create the second alert message
	message = NewAlertMessage(model.WithAllDependencies(ts.Dependencies), model.New())
	ts.Require().NotNil(message)
	ts.Require().Equal(uint64(0), message.GetID())

	message.Hash = testAlertHash + "2"
	message.Raw = testAlertRaw + "2"
	message.SequenceNumber = 2

	err = message.Save(context.Background())
	ts.Require().NoError(err)

	// Get the latest alert message
	message, err = GetLatestAlert(context.Background(), nil, model.WithAllDependencies(ts.Dependencies))
	ts.Require().NoError(err)

	ts.Require().Equal(uint64(2), message.GetID())
	ts.Require().Equal(testAlertHash+"2", message.Hash)
	ts.Require().Equal(testAlertRaw+"2", message.Raw)
	ts.Require().Equal(uint32(2), message.SequenceNumber)
}

// TestAlertMessage_SerializeData will test serializing the data
func (ts *TestSuite) TestAlertMessage_SerializeData() {
	message := NewAlertMessage(model.WithAllDependencies(ts.Dependencies), model.New())
	ts.Require().NotNil(message)
	ts.Require().Equal(uint64(0), message.GetID())

	message.Hash = testAlertHash
	message.Raw = testAlertRaw
	message.SequenceNumber = 1
	message.alertType = AlertTypeInformational

	err := message.Save(context.Background())
	ts.Require().NoError(err)

	message.SerializeData()

	ts.Require().Equal(uint64(1), message.GetID())
	ts.Equal("dea523f6449b08f64b8b4c1f416333bd14ed14ebe2d2585a826cd348228a3ecf", message.Hash)
	ts.Equal("0000000001000000000000000000000001000000", hex.EncodeToString(message.GetRawData()))
	ts.Equal(AlertTypeInformational, message.GetAlertType())
}
