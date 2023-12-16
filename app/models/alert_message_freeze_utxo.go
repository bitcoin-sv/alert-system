package models

import "context"

// AlertMessageFreezeUtxo is the message for freezing UTXOs
type AlertMessageFreezeUtxo struct {
	AlertMessage
	// TODO finish building out this alert type
}

// Read reads the message
func (a *AlertMessageFreezeUtxo) Read(_ []byte) error {
	return nil
}

// Do performs the message
func (a *AlertMessageFreezeUtxo) Do(_ context.Context) error {
	return nil
}
