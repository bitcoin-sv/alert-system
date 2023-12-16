// Package mocks is a generated mocking package for the mocks
package mocks

import "context"

// Node is a mock type for the SVNode interface
type Node struct {
	// Fields
	RPCHost     string
	RPCPassword string
	RPCUser     string

	// Functions
	BanPeerFunc         func(ctx context.Context, peer string) error
	InvalidateBlockFunc func(ctx context.Context, hash string) error
	UnbanPeerFunc       func(ctx context.Context, peer string) error

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
