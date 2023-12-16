package config

import (
	"context"

	"github.com/bitcoin-sv/alert-system/app/config/mocks"
	"github.com/libsv/go-bn"
)

// NodeInterface is the interface for a node
type NodeInterface interface {
	BanPeer(ctx context.Context, peer string) error
	GetRPCHost() string
	GetRPCPassword() string
	GetRPCUser() string
	InvalidateBlock(ctx context.Context, hash string) error
	UnbanPeer(ctx context.Context, peer string) error
}

// NewNodeConfig creates a new NodeConfig struct
func NewNodeConfig(user, pass, host string) NodeInterface {
	return &Node{
		RPCUser:     user,
		RPCPassword: pass,
		RPCHost:     host,
	}
}

// NewNodeMock creates a new NodeConfig struct for testing
func NewNodeMock(user, pass, host string) NodeInterface {
	return &mocks.Node{
		RPCUser:     user,
		RPCPassword: pass,
		RPCHost:     host,
	}
}

// GetRPCUser returns the RPC user
func (n *Node) GetRPCUser() string {
	return n.RPCUser
}

// GetRPCPassword returns the RPC password
func (n *Node) GetRPCPassword() string {
	return n.RPCPassword
}

// GetRPCHost returns the RPC host
func (n *Node) GetRPCHost() string {
	return n.RPCHost
}

// InvalidateBlock invalidates a block
func (n *Node) InvalidateBlock(ctx context.Context, hash string) error {
	c := bn.NewNodeClient(bn.WithCreds(n.RPCUser, n.RPCPassword), bn.WithHost(n.RPCHost))
	return c.InvalidateBlock(ctx, hash)
}

// BanPeer bans a peer
func (n *Node) BanPeer(ctx context.Context, peer string) error {
	c := bn.NewNodeClient(bn.WithCreds(n.RPCUser, n.RPCPassword), bn.WithHost(n.RPCHost))
	return c.SetBan(ctx, peer, bn.BanActionAdd, nil)
}

// UnbanPeer unbans a peer
func (n *Node) UnbanPeer(ctx context.Context, peer string) error {
	c := bn.NewNodeClient(bn.WithCreds(n.RPCUser, n.RPCPassword), bn.WithHost(n.RPCHost))
	return c.SetBan(ctx, peer, bn.BanActionRemove, nil)
}
