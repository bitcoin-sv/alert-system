package p2p

import (
	"context"
	"encoding/hex"
	"math"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libsv/go-p2p/wire"
)

// Thread is an interface for a thread
type Thread interface {
	Start(ctx context.Context) error
	Kill(ctx context.Context) error
}

// StreamThread is a thread for a stream
type StreamThread struct {
	config           *config.Config
	ctx              context.Context // TODO should remove this, should be passed in via methods only
	latestSequence   uint32
	myLatestSequence uint32
	peer             peer.ID
	stream           network.Stream
}

// LatestSequence will return the threads latest sequence
func (s *StreamThread) LatestSequence() uint32 {
	return s.latestSequence
}

// Sync will start the thread
func (s *StreamThread) Sync(ctx context.Context) error {

	// Get the latest alert
	a, err := models.GetLatestAlert(ctx, nil, model.WithAllDependencies(s.config))
	if err != nil {
		s.config.Services.Log.Errorf("failed to get latest alert: %s", err.Error())
		return err
	} else if a == nil {
		s.config.Services.Log.Error(ErrAlertNotLatest.Error())
		return ErrAlertNotLatest
	}

	s.myLatestSequence = a.SequenceNumber
	// construct get latest message
	msg := SyncMessage{
		Type: IWantLatest,
	}
	data := msg.Serialize()

	defer func() {
		_ = s.stream.Close()
	}()

	if err = wire.WriteVarBytes(s.stream, 0, data); err != nil {
		return err
	}

	s.config.Services.Log.Debugf("requested latest sequence in stream %s", s.stream.ID())

	return s.ProcessSyncMessage(ctx)
}

// ProcessSyncMessage will process the sync message
func (s *StreamThread) ProcessSyncMessage(ctx context.Context) error {
	for {
		b, err := wire.ReadVarBytes(s.stream, 0, math.MaxUint64, config.ApplicationName)
		if err != nil {
			if s.stream.Conn().IsClosed() || s.stream.Stat().Transient {
				return nil
			}
			s.config.Services.Log.Debugf("failed to read sync message: %s; closing stream", err.Error())
			return s.stream.Close()
		}

		if len(b) == 0 {
			_ = s.stream.Close()
			return nil
		}
		var msg *SyncMessage
		if msg, err = NewSyncMessageFromBytes(b); err != nil {
			s.config.Services.Log.Errorf("failed to convert to sync message: %s", err.Error())
			return err
		}
		switch msg.Type {
		case IGotLatest:
			s.config.Services.Log.Debugf("received latest sequence %d from peer %s", msg.SequenceNumber, s.peer.String())
			if err = s.ProcessGotLatest(ctx, msg); err != nil {
				return err
			}
			if s.myLatestSequence >= s.latestSequence {
				_ = s.stream.Close()
				return nil
			}
			s.config.Services.Log.Debugf("wrote msg requesting next sequence %d from peer %s", s.myLatestSequence+1, s.peer.String())
		case IGotSequenceNumber:
			s.config.Services.Log.Debugf("received IGotSequenceNumber %d from peer %s", msg.SequenceNumber, s.peer.String())
			if err = s.ProcessGotSequenceNumber(msg); err != nil {
				return err
			}
			if s.myLatestSequence == s.latestSequence {
				_ = s.stream.Close()
				return nil
			}
			s.config.Services.Log.Debugf("wrote msg requesting next sequence %d from peer %s", msg.SequenceNumber+1, s.peer.String())
		case IWantSequenceNumber:
			s.config.Services.Log.Debugf("received IWantSequenceNumber %d from peer %s", msg.SequenceNumber, s.peer.String())
			if err = s.ProcessWantSequenceNumber(ctx, msg); err != nil {
				return err
			}
			s.config.Services.Log.Debugf("wrote sequence %d to peer %s", msg.SequenceNumber, s.peer.String())
			if msg.SequenceNumber == s.myLatestSequence {
				_ = s.stream.Close()
				return nil
			}
		case IWantLatest:
			s.config.Services.Log.Debugf("received IWantLatest from peer %s", s.peer.String())
			if err = s.ProcessWantLatest(ctx); err != nil {
				return err
			}
			s.config.Services.Log.Debugf("wrote latest sequence %d to peer %s", s.myLatestSequence, s.peer.String())
		}
	}
}

// ProcessGotLatest will process the got latest message
func (s *StreamThread) ProcessGotLatest(ctx context.Context, msg *SyncMessage) error {
	a, err := models.GetLatestAlert(ctx, nil, model.WithAllDependencies(s.config))
	if err != nil {
		s.config.Services.Log.Errorf("failed to get latest alert to send to peer: %s", err.Error())
		return err
	} else if a == nil {
		s.config.Services.Log.Error(ErrAlertNotLatest.Error())
		return ErrAlertNotLatest
	}

	s.myLatestSequence = a.SequenceNumber // this is redundant, but doesn't hurt
	if msg.SequenceNumber < a.SequenceNumber {
		s.config.Services.Log.Debugf("peer %s is not synced yet, ignoring...", s.peer.String())
		return nil
	}

	s.latestSequence = msg.SequenceNumber
	if msg.SequenceNumber == a.SequenceNumber {
		s.config.Services.Log.Debugf("peer %s is synced to current state as us, closing stream.", s.peer.String())
		_ = s.stream.Close()
		return nil
	}
	s.config.Services.Log.Infof("peer %s has sequence %d and we have %d", s.peer.String(), msg.SequenceNumber, a.SequenceNumber)

	// need to get next sequence
	res := SyncMessage{
		Type:           IWantSequenceNumber,
		SequenceNumber: a.SequenceNumber + 1,
	}
	return wire.WriteVarBytes(s.stream, 0, res.Serialize())
}

// ProcessGotSequenceNumber will process the got sequence number message
func (s *StreamThread) ProcessGotSequenceNumber(msg *SyncMessage) error {
	// Sync with a new alert
	a, err := models.NewAlertFromBytes(msg.Data, model.WithAllDependencies(s.config), model.New())
	if err != nil {
		// todo probably want to ban this peer?
		return err
	}

	// Verify signatures
	var valid bool
	if valid, err = a.AreSignaturesValid(s.ctx); err != nil {
		return err
	} else if !valid { // Not valid
		s.config.Services.Log.Error(ErrInvalidAlerts.Error())
		return ErrInvalidAlerts
	}

	// Serialize the alert data and hash
	a.SerializeData()

	// Process the alert (if it's a set keys alert)
	if a.GetAlertType() == models.AlertTypeSetKeys {
		ak := a.ProcessAlertMessage()
		if err = ak.Read(a.GetRawMessage()); err != nil {
			return err
		}
		if err = ak.Do(s.ctx); err != nil {
			return err
		}
	}

	// Save the alert
	if err = a.Save(s.ctx); err != nil {
		return err
	}

	// Update the latest sequence
	s.myLatestSequence = a.SequenceNumber
	if s.myLatestSequence == s.latestSequence {
		s.config.Services.Log.Infof("successfully synced up to sequence %d", s.latestSequence)
		_ = s.stream.Close()
		return nil
	}

	// need to get next sequence
	res := SyncMessage{
		Type:           IWantSequenceNumber,
		SequenceNumber: a.SequenceNumber + 1,
	}
	return wire.WriteVarBytes(s.stream, 0, res.Serialize())
}

// ProcessWantSequenceNumber will process the want sequence number message
func (s *StreamThread) ProcessWantSequenceNumber(ctx context.Context, msg *SyncMessage) error {
	a, err := models.GetAlertMessageBySequenceNumber(ctx, msg.SequenceNumber, model.WithAllDependencies(s.config))
	if err != nil {
		s.config.Services.Log.Errorf("failed to get latest alert to send to peer: %s", err.Error())
		return err
	} else if a == nil {
		s.config.Services.Log.Error(ErrAlertNotFoundBySequence.Error())
		return ErrAlertNotFoundBySequence
	}
	var data []byte
	if data, err = hex.DecodeString(a.Raw); err != nil {
		s.config.Services.Log.Errorf("failed to decode raw alert data: %s", err.Error())
		return err
	}
	res := SyncMessage{
		Type:           IGotSequenceNumber,
		SequenceNumber: a.SequenceNumber,
		Data:           data,
	}
	return wire.WriteVarBytes(s.stream, 0, res.Serialize())
}

// ProcessWantLatest will process the want latest message
func (s *StreamThread) ProcessWantLatest(ctx context.Context) error {
	a, err := models.GetLatestAlert(ctx, nil, model.WithAllDependencies(s.config))
	if err != nil {
		s.config.Services.Log.Errorf("failed to get latest alert to send to peer: %s", err.Error())
		return err
	} else if a == nil {
		s.config.Services.Log.Error(ErrAlertNotLatest.Error())
		return ErrAlertNotLatest
	}
	s.myLatestSequence = a.SequenceNumber

	var data []byte
	if data, err = hex.DecodeString(a.Raw); err != nil {
		s.config.Services.Log.Errorf("failed to decode raw alert data: %s", err.Error())
		return err
	}
	res := SyncMessage{
		Type:           IGotLatest,
		SequenceNumber: a.SequenceNumber,
		Data:           data,
	}
	return wire.WriteVarBytes(s.stream, 0, res.Serialize())
}
