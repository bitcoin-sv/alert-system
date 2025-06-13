package p2p

import (
	"encoding/binary"
)

// IWantLatest is the byte for "I want the latest"
const IWantLatest = 0x01

// IWantSequenceNumber is the byte for "I want sequence number"
const IWantSequenceNumber = 0x02

// IGotSequenceNumber is the byte for "I got sequence number"
const IGotSequenceNumber = 0x03

// IGotLatest is the byte for "I got latest"
const IGotLatest = 0x04

// SyncMessage is the message for syncing
type SyncMessage struct {
	Data           []byte `json:"data"`
	SequenceNumber uint32 `json:"sequence_number"`
	Type           byte   `json:"type"`
}

// NewSyncMessageFromBytes will create a new sync message from bytes
func NewSyncMessageFromBytes(in []byte) (*SyncMessage, error) {
	if len(in) < 1 {
		return nil, ErrSyncMessageByte
	}
	s := SyncMessage{}
	s.Type = in[0]
	if s.Type == IWantLatest {
		return &s, nil
	}
	if len(in) < 5 {
		return nil, ErrSyncFiveBytes
	}
	s.SequenceNumber = binary.LittleEndian.Uint32(in[1:5])
	s.Data = in[5:]
	return &s, nil
}

// Serialize will serialize the sync message
func (s *SyncMessage) Serialize() []byte {
	var ret []byte
	ret = append(ret, s.Type)
	ret = binary.LittleEndian.AppendUint32(ret, s.SequenceNumber)
	ret = append(ret, s.Data...)
	return ret
}
