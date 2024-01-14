package models

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/libsv/go-bn/models"
)

// AlertMessageConfiscateTransaction is a confiscate utxo alert
type AlertMessageConfiscateTransaction struct {
	AlertMessage
	Transactions []models.ConfiscationTransactionDetails
}

type ConfiscateTransaction struct {
	EnforceAtHeight [8]byte
	Id              [32]byte
}

// Read reads the alert
func (a *AlertMessageConfiscateTransaction) Read(raw []byte) error {
	a.Config().Services.Log.Infof("%x", raw)
	if len(raw) < 40 {
		return fmt.Errorf("confiscation alert is less than 41 bytes")
	}
	if len(raw)%40 != 0 {
		return fmt.Errorf("confiscation alert is not a multiple of 41 bytes")
	}
	txCount := len(raw) / 40
	details := []models.ConfiscationTransactionDetails{}
	for i := 0; i < txCount; i++ {
		tx := ConfiscateTransaction{
			EnforceAtHeight: [8]byte(raw[:8]),
			Id:              [32]byte(raw[8:40]),
		}
		detail := models.ConfiscationTransactionDetails{
			ConfiscationTransaction: models.ConfiscationTransaction{
				EnforceAtHeight: int64(binary.LittleEndian.Uint64(tx.EnforceAtHeight[:])),
				Hex:             hex.EncodeToString(tx.Id[:]),
			},
		}
		details = append(details, detail)
		raw = raw[40:]
	}
	return nil
}

// Do executes the alert
func (a *AlertMessageConfiscateTransaction) Do(ctx context.Context) error {
	_, err := a.Config().Services.Node.AddToConfiscationTransactionWhitelist(ctx, a.Transactions)
	if err != nil {
		return err
	}
	return nil
}
