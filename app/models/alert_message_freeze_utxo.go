package models

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/bsv-blockchain/go-bn/models"
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
	var raw []byte
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
	var funds []models.Fund
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

// Do perform the message
func (a *AlertMessageFreezeUtxo) Do(ctx context.Context) error {
	_, err := a.Config().Services.Node.AddToConsensusBlacklist(ctx, a.Funds)
	if err != nil {
		return err
	}
	return nil
}

// ToJSON is the alert in JSON format
func (a *AlertMessageFreezeUtxo) ToJSON(_ context.Context) []byte {
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
func (a *AlertMessageFreezeUtxo) MessageString() string {
	return fmt.Sprintf("Freezing utxo id [%x]; vout: [%d], enforcing at height start [%d], end [%d].", a.Funds[0].TxOut.TxId, a.Funds[0].TxOut.Vout, a.Funds[0].EnforceAtHeight[0].Start, a.Funds[0].EnforceAtHeight[0].Stop)
}
