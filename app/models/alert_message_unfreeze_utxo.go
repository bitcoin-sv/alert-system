package models

import "context"

// AlertMessageUnfreezeUtxo is the message for unfreezing a UTXO
type AlertMessageUnfreezeUtxo struct {
	AlertMessage
	// TODO finish building out this alert type
	Funds []Fund
}

// Read reads the message from the byte slice
func (a *AlertMessageUnfreezeUtxo) Read(_ []byte) error {
	return nil
}

// Do executes the message
func (a *AlertMessageUnfreezeUtxo) Do(_ context.Context) error {
	return nil
}
