package models

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/libsv/go-bn/models"
)

// AlertMessageUnfreezeUtxo is the message for unfreezing a UTXO
type AlertMessageUnfreezeUtxo struct {
	AlertMessage
	// TODO finish building out this alert type
	Funds []models.Fund
}

// Read reads the message from the byte slice
func (a *AlertMessageUnfreezeUtxo) Read(raw []byte) error {
	if len(raw) < 57 {
		return fmt.Errorf("unfreeze alert is less than 57 bytes, got %d bytes; raw: %x", len(raw), raw)
	}
	if len(raw)%57 != 0 {
		return fmt.Errorf("unfreeze alert is not a multiple of 57 bytes, got %d bytes; raw: %x", len(raw), raw)
	}
	fundCount := len(raw) / 57
	funds := []models.Fund{}
	for i := 0; i < fundCount; i++ {
		fund := Fund{
			TransactionOutID:     [32]byte(raw[0:32]),
			Vout:                 binary.LittleEndian.Uint64(raw[32:40]),
			EnforceAtHeightStart: binary.LittleEndian.Uint64(raw[40:48]),
			EnforceAtHeightEnd:   binary.LittleEndian.Uint64(raw[48:56]),
		}
		enforceByte := raw[56]

		if enforceByte != uint8(0) {
			fund.PolicyExpiresWithConsensus = true
		}
		funds = append(funds, models.Fund{
			TxOut: models.TxOut{
				TxId: hex.EncodeToString(fund.TransactionOutID[:]),
				Vout: int(fund.Vout),
			},
			EnforceAtHeight: []models.Enforce{
				{
					Start: int(fund.EnforceAtHeightStart),
					Stop:  int(fund.EnforceAtHeightEnd),
				},
			},
			PolicyExpiresWithConsensus: fund.PolicyExpiresWithConsensus,
		})
		raw = raw[57:]
	}
	a.Funds = funds

	return nil

}

// Do executes the message
func (a *AlertMessageUnfreezeUtxo) Do(ctx context.Context) error {
	_, err := a.Config().Services.Node.AddToConsensusBlacklist(ctx, a.Funds)
	if err != nil {
		return err
	}
	return nil
}
