package models

import (
	"bytes"
	"context"
	"fmt"

	"github.com/libsv/go-p2p/wire"
)

// AlertMessageUnbanPeer is the message for unban peer
type AlertMessageUnbanPeer struct {
	AlertMessage
	Peer         []byte `json:"peer"`
	PeerLength   uint64 `json:"peer_length"`
	Reason       []byte `json:"reason"`
	ReasonLength uint64 `json:"reason_length"`
}

// Read reads the payload from the byte slice
func (a *AlertMessageUnbanPeer) Read(alert []byte) error {
	buf := bytes.NewReader(alert)

	// read the peer length
	peerLength, err := wire.ReadVarInt(buf, 0)
	if err != nil {
		return err
	}

	// read the peer IP + port
	var peer []byte
	for i := uint64(0); i < peerLength; i++ {
		var b byte
		if b, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("failed to read peer: %s", err.Error())
		}
		peer = append(peer, b)
	}
	a.PeerLength = peerLength
	a.Peer = peer

	// read the reason
	var reasonLength uint64
	if reasonLength, err = wire.ReadVarInt(buf, 0); err != nil {
		return err
	}
	var reason []byte
	for i := uint64(0); i < reasonLength; i++ {
		var b byte
		if b, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("failed to read reason: %s", err.Error())
		}
		reason = append(reason, b)
	}

	a.Reason = reason
	a.ReasonLength = reasonLength
	return nil
}

// Do executes the alert
func (a *AlertMessageUnbanPeer) Do(ctx context.Context) error {
	return a.Config().Services.Node.UnbanPeer(ctx, string(a.Peer))
}
