package models

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/libsv/go-bn/models"
)

// AlertMessageFreezeUtxo is the message for freezing UTXOs
type AlertMessageFreezeUtxo struct {
	AlertMessage
	Funds []models.Fund
}

// Fund is the struct defining funds to freeze
type Fund struct {
	TransactionOutID           [32]byte
	Vout                       uint64
	EnforceAtHeightStart       uint64
	EnforceAtHeightEnd         uint64
	PolicyExpiresWithConsensus bool
}

// Serialize creates the raw hex string of the fund
func (f *Fund) Serialize() []byte {
	raw := []byte{}
	raw = append(raw, f.TransactionOutID[:]...)
	raw = binary.LittleEndian.AppendUint64(raw, f.Vout)
	raw = binary.LittleEndian.AppendUint64(raw, f.EnforceAtHeightStart)
	raw = binary.LittleEndian.AppendUint64(raw, f.EnforceAtHeightEnd)
	expire := uint8(0)
	if f.PolicyExpiresWithConsensus {
		expire = uint8(1)
	}
	raw = append(raw, expire)
	return raw
}

// Read reads the message
func (a *AlertMessageFreezeUtxo) Read(raw []byte) error {
	if len(raw) < 57 {
		return fmt.Errorf("freeze alert is less than 57 bytes, got %d bytes; raw: %x", len(raw), raw)
	}
	if len(raw)%57 != 0 {
		return fmt.Errorf("freeze alert is not a multiple of 57 bytes, got %d bytes; raw: %x", len(raw), raw)
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

// Do performs the message
func (a *AlertMessageFreezeUtxo) Do(ctx context.Context) error {
	_, err := a.Config().Services.Node.AddToConsensusBlacklist(ctx, a.Funds)
	if err != nil {
		return err
	}
	return nil
}
