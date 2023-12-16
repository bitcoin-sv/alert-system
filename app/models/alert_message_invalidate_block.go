package models

import (
	"bytes"
	"context"
	"fmt"

	"github.com/libsv/go-bt/v2/chainhash"
	"github.com/libsv/go-p2p/wire"
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

	buf := bytes.NewReader(alert[32:])

	// read the reason length
	var length uint64
	if length, err = wire.ReadVarInt(buf, 0); err != nil {
		return err
	}
	var msg []byte
	for i := uint64(0); i < length; i++ {
		var b byte
		if b, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("failed to read reason: %s", err.Error())
		}
		msg = append(msg, b)
	}
	a.ReasonLength = length
	a.Reason = msg
	a.BlockHash = blockHash
	return nil
}

// Do executes the alert
func (a *AlertMessageInvalidateBlock) Do(ctx context.Context) error {
	return a.Config().Services.Node.InvalidateBlock(ctx, a.BlockHash.String())
}
