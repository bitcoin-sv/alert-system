package models

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/libsv/go-bn/models"
	"github.com/libsv/go-p2p/wire"
)

// AlertMessageConfiscateTransaction is a confiscate utxo alert
type AlertMessageConfiscateTransaction struct {
	AlertMessage
	Transactions []models.ConfiscationTransactionDetails
}

// ConfiscateTransaction defines the parameters for the confiscation transaction
type ConfiscateTransaction struct {
	EnforceAtHeight uint64
	Hex             []byte
}

// Read reads the alert
func (a *AlertMessageConfiscateTransaction) Read(raw []byte) error {
	a.Config().Services.Log.Infof("%x", raw)
	if len(raw) < 9 {
		return fmt.Errorf("confiscation alert is less than 9 bytes")
	}
	// TODO: assume for now only 1 confiscation tx in the alert for simplicity
	var details []models.ConfiscationTransactionDetails
	enforceAtHeight := binary.LittleEndian.Uint64(raw[0:8])
	buf := bytes.NewReader(raw[8:])

	length, err := wire.ReadVarInt(buf, 0)
	if err != nil {
		return err
	}
	if length > uint64(buf.Len()) {
		return errors.New("tx hex length is longer than the remaining buffer")
	}

	// read the tx hex
	var rawHex []byte
	for i := uint64(0); i < length; i++ {
		var b byte
		if b, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("failed to read tx hex: %s", err.Error())
		}
		rawHex = append(rawHex, b)
	}

	detail := models.ConfiscationTransactionDetails{
		ConfiscationTransaction: models.ConfiscationTransaction{
			EnforceAtHeight: int64(enforceAtHeight),
			Hex:             hex.EncodeToString(rawHex),
		},
	}
	details = append(details, detail)

	a.Transactions = details
	a.Config().Services.Log.Infof("ConfiscateTransaction alert; enforceAt [%d]; hex [%s]", enforceAtHeight, hex.EncodeToString(rawHex))

	return nil
}

// Do execute the alert
func (a *AlertMessageConfiscateTransaction) Do(ctx context.Context) error {
	res, err := a.Config().Services.Node.AddToConfiscationTransactionWhitelist(ctx, a.Transactions)
	if err != nil {
		return err
	}
	if len(res.NotProcessed) > 0 {
		// we can safely assume this is just one not processed tx because we are only publishing one tx with the alert right now
		return fmt.Errorf("confiscation alert RPC response returned an error; reason: %s", res.NotProcessed[0].Reason)
	}
	return nil
}

// ToJSON is the alert in JSON format
func (a *AlertMessageConfiscateTransaction) ToJSON(_ context.Context) []byte {
	m := a.ProcessAlertMessage()
	// TODO: Come back and add a message interface for each alert
	_ = m.Read(a.GetRawMessage())
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return []byte{}
	}
	return data
}

// MessageString executes the alert
func (a *AlertMessageConfiscateTransaction) MessageString() string {
	return fmt.Sprintf("Adding confiscation transaction [%x] to whitelist enforcing at height [%d].", a.Transactions[0].ConfiscationTransaction.Hex, a.Transactions[0].ConfiscationTransaction.EnforceAtHeight)
}
