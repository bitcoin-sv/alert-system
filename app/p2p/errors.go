package p2p

import "errors"

// Errors for the p2p package
var (
	ErrAlertNotFoundBySequence = errors.New("failed to find alert by sequence in datastore")
	ErrAlertNotLatest          = errors.New("failed to find latest alert datastore")
	ErrInvalidAlerts           = errors.New("peer is sending invalid alerts")
	ErrSyncFiveBytes           = errors.New("sync message is less than 5 bytes, not valid")
	ErrSyncMessageByte         = errors.New("sync message needs at least a byte")
)
