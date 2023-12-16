package models

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/alert-system/app/models/model"
)

// AlertMessageSetKeys is the message for setting keys
type AlertMessageSetKeys struct {
	AlertMessage
	Keys [][33]byte
	Hash string
}

// Read reads the message
func (a *AlertMessageSetKeys) Read(alert []byte) error {

	// Check the length
	if len(alert) != 165 {
		return fmt.Errorf("alert is not 165 bytes long, got %d bytes, not valid", len(alert))
	}
	buf := bytes.NewReader(alert[:])

	// Read the message hash
	for key := 0; key < 5; key++ {
		var pubKey []byte
		for i := uint64(0); i < 33; i++ {
			b, err := buf.ReadByte()
			if err != nil {
				return fmt.Errorf("failed to read pubKey: %s", err.Error())
			}
			pubKey = append(pubKey, b)
		}
		a.Keys = append(a.Keys, [33]byte(pubKey))
	}

	return nil
}

// Do executes the alert
func (a *AlertMessageSetKeys) Do(ctx context.Context) error {
	err := ClearActivePublicKeys(ctx, a.Config().Services.Datastore)
	if err != nil {
		return err
	}
	for _, key := range a.Keys {
		pk := NewPublicKey(model.WithAllDependencies(a.Config()))
		pk.Key = hex.EncodeToString(key[:])
		pk.Active = true
		pk.LastUpdateHash = a.Hash
		if err = pk.Save(ctx); err != nil {
			return err
		}
	}
	return nil
}
