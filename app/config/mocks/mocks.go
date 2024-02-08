// Package mocks is a generated mocking package for the mocks
package mocks

import (
	"context"

	"github.com/libsv/go-bn/models"
)

// Node is a mock type for the SVNode interface
type Node struct {
	// Fields
	RPCHost     string
	RPCPassword string
	RPCUser     string

	// Functions
	BanPeerFunc                               func(ctx context.Context, peer string) error
	BestBlockHashFunc                         func(ctx context.Context) (string, error)
	InvalidateBlockFunc                       func(ctx context.Context, hash string) error
	UnbanPeerFunc                             func(ctx context.Context, peer string) error
	AddToConsensusBlacklistFunc               func(ctx context.Context, funds []models.Fund) (*models.AddToConsensusBlacklistResponse, error)
	AddToConfiscationTransactionWhitelistFunc func(ctx context.Context, tx []models.ConfiscationTransactionDetails) (*models.AddToConfiscationTransactionWhitelistResponse, error)
	// Add additional fields if needed to track calls or results
}

// GetRPCUser will return the RPCUser
func (n *Node) GetRPCUser() string {
	return n.RPCUser
}

// GetRPCPassword will  return the RPCPassword
func (n *Node) GetRPCPassword() string {
	return n.RPCPassword
}

// GetRPCHost will return the RPCHost
func (n *Node) GetRPCHost() string {
	return n.RPCHost
}

// BanPeer will call the BanPeerFunc if not nil, otherwise return nil
func (n *Node) BanPeer(ctx context.Context, peer string) error {
	if n.BanPeerFunc != nil {
		return n.BanPeerFunc(ctx, peer)
	}
	// Default behavior if no mock function provided
	return nil
}

// BestBlockHash will call the BestBlockHashFunc
func (n *Node) BestBlockHash(ctx context.Context) (string, error) {
	if n.BestBlockHashFunc != nil {
		return n.BestBlockHashFunc(ctx)
	}
	return "", nil
}

// InvalidateBlock will call the InvalidateBlockFunc if not nil, otherwise return nil
func (n *Node) InvalidateBlock(ctx context.Context, hash string) error {
	if n.InvalidateBlockFunc != nil {
		return n.InvalidateBlockFunc(ctx, hash)
	}
	return nil
}

// UnbanPeer will call the UnbanPeerFunc if not nil, otherwise return nil
func (n *Node) UnbanPeer(ctx context.Context, peer string) error {
	if n.UnbanPeerFunc != nil {
		return n.UnbanPeerFunc(ctx, peer)
	}
	return nil
}

// AddToConsensusBlacklist will call the AddToConsensusBlacklistFunc if not nil, otherwise return nil
func (n *Node) AddToConsensusBlacklist(ctx context.Context, funds []models.Fund) (*models.AddToConsensusBlacklistResponse, error) {
	if n.AddToConsensusBlacklistFunc != nil {
		return n.AddToConsensusBlacklistFunc(ctx, funds)
	}
	return nil, nil
}

// AddToConfiscationTransactionWhitelist will call the AddToConfiscationTransactionWhitelistFunc if not nil, otherwise return nil
func (n *Node) AddToConfiscationTransactionWhitelist(ctx context.Context, tx []models.ConfiscationTransactionDetails) (*models.AddToConfiscationTransactionWhitelistResponse, error) {
	if n.AddToConfiscationTransactionWhitelistFunc != nil {
		return n.AddToConfiscationTransactionWhitelistFunc(ctx, tx)
	}
	return nil, nil
}
