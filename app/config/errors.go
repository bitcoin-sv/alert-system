package config

import (
	"errors"
)

// Configuration errors
var (
	ErrDatastoreRequired    = errors.New("datastore is required and was not loaded")
	ErrDatastoreUnsupported = errors.New("unsupported datastore engine")
	ErrInvalidEnvironment   = errors.New("invalid environment")
	ErrNoP2PIP              = errors.New("no p2p_ip defined")
	ErrNoP2PPort            = errors.New("no p2p_port defined")
	ErrNoRPCHost            = errors.New("no rpc_host defined")
	ErrNoRPCPassword        = errors.New("no rpc_password defined")
	ErrNoRPCUser            = errors.New("no rpc_user defined")
	ErrNoRPCConnections     = errors.New("no rpc connections configured")
)
