package models

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
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
	details := []models.ConfiscationTransactionDetails{}
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

// Do executes the alert
func (a *AlertMessageConfiscateTransaction) Do(ctx context.Context) error {
	res, err := a.Config().Services.Node.AddToConfiscationTransactionWhitelist(ctx, a.Transactions)
	if err != nil {
		return err
	}
	if len(res.NotProcessed) > 0 {
		a.Config().Services.Log.Errorf("confiscation alert RPC response indicates it might have not been processed")
		// TODO: I think we want to error here in the future so that the RPC call will be retried... but not clear right now
	}
	return nil
}
