package models

import (
	"context"
	"time"

	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/bitcoin-sv/alert-system/utils"
)

// CreateGenesisAlert will create the genesis alert if it is not in the database
// If it is in the database, it will do nothing
func CreateGenesisAlert(ctx context.Context, opts ...model.Options) error {
	newAlert := NewAlertMessage(opts...)
	// Get the alert message by sequence number
	m, err := GetAlertMessageBySequenceNumber(ctx, 0, opts...)
	if err != nil {
		return err
	} else if m != nil && len(m.Hash) > 0 { // found it, skipping (no need to create)
		return nil
	}

	// Create the array of keys
	keys := newAlert.Config().GenesisKeys
	var msg []byte

	// Create the array of keys to save
	keysToSave := make([]*PublicKey, 0, len(keys))

	// Set the new flag
	opts = append(opts, model.New())

	// Loop through the keys
	for _, key := range keys {
		// Create the public key
		k := NewPublicKey(opts...)
		k.Key = key
		k.Active = true
		keysToSave = append(keysToSave, k)
	}

	// Sync creating a new alert

	newAlert.SetAlertType(AlertTypeSetKeys)
	newAlert.message = msg
	newAlert.SequenceNumber = 0
	newAlert.timestamp = uint64(time.Date(2923, time.November, 1, 1, 1, 1, 1, time.UTC).Unix())
	newAlert.version = 1
	newAlert.Processed = true

	// Serialize the data
	newAlert.SerializeData()

	// Sign the genesis alert
	var sigs [][]byte
	if sigs, err = utils.SignWithGenesis(newAlert.data); err != nil {
		return err
	}
	newAlert.signatures = sigs

	// Save the keys
	for i := range keysToSave {
		if err = keysToSave[i].Save(ctx); err != nil {
			return err
		}
	}
	_ = newAlert.Serialize()

	// Save the alert
	return newAlert.Save(ctx)
}
