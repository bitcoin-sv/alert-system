package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/libsv/go-p2p/wire"
)

// AlertMessageBanPeer is the message for ban peer
type AlertMessageBanPeer struct {
	AlertMessage
	Peer         []byte `json:"peer"`
	PeerLength   uint64 `json:"peer_length"`
	Reason       []byte `json:"reason"`
	ReasonLength uint64 `json:"reason_length"`
}

// Read reads the payload from the byte slice
func (a *AlertMessageBanPeer) Read(alert []byte) error {
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
func (a *AlertMessageBanPeer) Do(ctx context.Context) error {
	return a.Config().Services.Node.BanPeer(ctx, string(a.Peer))
}

// ToJSON is the alert in JSON format
func (a *AlertMessageBanPeer) ToJSON(_ context.Context) []byte {
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
func (a *AlertMessageBanPeer) MessageString() string {
	return fmt.Sprintf("Banning peer [%s]; reason [%s].", a.Peer, a.Reason)
}
