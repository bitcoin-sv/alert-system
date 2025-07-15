package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bsv-blockchain/go-sdk/util"

	"github.com/bsv-blockchain/go-bt/v2/chainhash"
)

// AlertMessageInvalidateBlock is an invalidate block alert
type AlertMessageInvalidateBlock struct {
	AlertMessage
	BlockHash    *chainhash.Hash `json:"block_hash"`
	ReasonLength uint64          `json:"reason_length"`
	Reason       []byte          `json:"reason"`
}

// Read reads the alert
func (a *AlertMessageInvalidateBlock) Read(alert []byte) error {
	blockHash, err := chainhash.NewHash(alert[:32])
	if err != nil {
		return err
	}

	reader := util.NewReader(alert[32:])

	// read the reason length
	var length uint64
	if length, err = reader.ReadVarInt(); err != nil {
		return err
	}
	if length == 0 {
		return fmt.Errorf("no reason message provided")
	}
	var msg []byte
	for i := uint64(0); i < length; i++ {
		var b byte
		if b, err = reader.ReadByte(); err != nil {
			return fmt.Errorf("failed to read reason: %s", err.Error())
		}
		msg = append(msg, b)
	}
	if !reader.IsComplete() {
		return fmt.Errorf("too many bytes in alert message")
	}
	a.ReasonLength = length
	a.Reason = msg
	a.BlockHash = blockHash
	return nil
}

// Do execute the alert
func (a *AlertMessageInvalidateBlock) Do(ctx context.Context) error {
	a.Config().Services.Log.Infof("InvalidateBlock alert; hash [%s]; reason [%s]", a.BlockHash, a.Reason)
	return a.Config().Services.Node.InvalidateBlock(ctx, a.BlockHash.String())
}

// ToJSON is the alert in JSON format
func (a *AlertMessageInvalidateBlock) ToJSON(_ context.Context) []byte {
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
func (a *AlertMessageInvalidateBlock) MessageString() string {
	return fmt.Sprintf("Invalidating block hash [%s]; reason [%s].", a.BlockHash, a.Reason)
}
