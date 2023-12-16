package models

import "context"

// AlertMessageConfiscateUtxo is a confiscate utxo alert
type AlertMessageConfiscateUtxo struct {
	AlertMessage
	// TODO finish building out this alert type
}

// Read reads the alert
func (a *AlertMessageConfiscateUtxo) Read(_ []byte) error {
	return nil
}

// Do executes the alert
func (a *AlertMessageConfiscateUtxo) Do(_ context.Context) error {
	return nil
}
