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

type Fund struct {
	TransactionOutId           [32]byte
	Vout                       [8]byte
	EnforceAtHeightStart       [8]byte
	EnforceAtHeightEnd         [8]byte
	PolicyExpiresWithConsensus bool
}

// Read reads the message
func (a *AlertMessageFreezeUtxo) Read(raw []byte) error {
	if len(raw) < 57 {
		return fmt.Errorf("freeze alert is less than 58 bytes")
	}
	if len(raw)%57 != 0 {
		return fmt.Errorf("freeze alert is not a multiple of 58 bytes")
	}
	fundCount := len(raw) / 57
	funds := []models.Fund{}
	for i := 0; i < fundCount; i++ {
		fund := Fund{
			TransactionOutId:     [32]byte(raw[0:32]),
			Vout:                 [8]byte(raw[32:40]),
			EnforceAtHeightStart: [8]byte(raw[40:48]),
			EnforceAtHeightEnd:   [8]byte(raw[48:56]),
		}
		enforceByte := binary.LittleEndian.Uint16(raw[56:57])

		if enforceByte != uint16(0) {
			fund.PolicyExpiresWithConsensus = true
		}
		funds = append(funds, models.Fund{
			TxOut: models.TxOut{
				TxId: hex.EncodeToString(fund.TransactionOutId[:]),
				Vout: int(binary.LittleEndian.Uint64(fund.Vout[:])),
			},
			EnforceAtHeight: []models.Enforce{
				{
					Start: int(binary.LittleEndian.Uint64(fund.EnforceAtHeightStart[:])),
					Stop:  int(binary.LittleEndian.Uint64(fund.EnforceAtHeightEnd[:])),
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
