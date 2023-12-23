package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/bitcoin-sv/alert-system/app/config/mocks"
	"github.com/stretchr/testify/require"
)

// TestBanPeer tests the BanPeer method
func TestBanPeer(t *testing.T) {
	mockNode := &mocks.Node{
		BanPeerFunc: func(ctx context.Context, peer string) error {
			// Mock behavior here
			if peer == "expected_peer_address" {
				return nil
			}
			return fmt.Errorf("unexpected peer address")
		},
	}

	ctx := context.Background()
	err := mockNode.BanPeer(ctx, "expected_peer_address")
	require.NoError(t, err)
}

// TestUnBanPeer tests the UnBanPeer method
func TestUnBanPeer(t *testing.T) {
	mockNode := &mocks.Node{
		UnbanPeerFunc: func(ctx context.Context, peer string) error {
			// Mock behavior here
			if peer == "expected_peer_address" {
				return nil
			}
			return fmt.Errorf("unexpected peer address")
		},
	}

	ctx := context.Background()
	err := mockNode.UnbanPeer(ctx, "expected_peer_address")
	require.NoError(t, err)
}

// TestInvalidateBlock tests the InvalidateBlock method
func TestInvalidateBlock(t *testing.T) {
	mockNode := &mocks.Node{
		InvalidateBlockFunc: func(ctx context.Context, hash string) error {
			// Mock behavior here
			if hash == "expected_hash" {
				return nil
			}
			return fmt.Errorf("unexpected hash")
		},
	}

	ctx := context.Background()
	err := mockNode.InvalidateBlock(ctx, "expected_hash")
	require.NoError(t, err)
}
