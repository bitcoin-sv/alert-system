package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/bitcoin-sv/alert-system/app/config/mocks"
	"github.com/bitcoin-sv/alert-system/app/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig_Failure tests the success case of LoadConfig
func TestLoadConfig_Success(t *testing.T) {
	// Setup
	tester.SetupEnv(t)
	defer func() {
		tester.TeardownEnv(t)
	}()

	// Execute
	config, err := LoadConfig(context.Background(), nil, true)
	require.NoError(t, err)

	// Assert
	assert.NotNil(t, config)
	assert.Equal(t, "user", config.RPCUser)
	assert.Equal(t, "password", config.RPCPassword)
	assert.Equal(t, "localhost", config.RPCHost)
	assert.Equal(t, "/path/to/private/key", config.P2PPrivateKeyPath)
	assert.Equal(t, SeedIpfsNode, config.P2PBootstrapPeer)
	assert.Equal(t, DefaultAlertSystemProtocolID, config.P2PAlertSystemProtocolID)
	assert.Equal(t, "192.168.1.1", config.P2PIP)
	assert.Equal(t, "8000", config.P2PPort)
	assert.Equal(t, "https://webhook.url", config.AlertWebhookURL)
}

// TestLoadConfig_MissingRPCUser tests the failure case of LoadConfig when RPC_USER is missing
func TestLoadConfig_MissingRPCUser(t *testing.T) {
	// Setup
	tester.SetupEnv(t)
	tester.UnsetEnv(t, "RPC_USER")

	defer func() {
		tester.TeardownEnv(t)
	}()

	// Execute
	_, err := LoadConfig(context.Background(), nil, true)

	// Assert
	require.Error(t, err)
	assert.Equal(t, ErrNoRPCUser, err)
}

// TestLoadConfig_MissingRPCPassword tests the failure case of LoadConfig when RPC_PASSWORD is missing
func TestLoadConfig_MissingRPCPassword(t *testing.T) {
	// Setup
	tester.SetupEnv(t)
	tester.UnsetEnv(t, "RPC_PASSWORD")

	defer func() {
		tester.TeardownEnv(t)
	}()

	// Execute
	_, err := LoadConfig(context.Background(), nil, true)

	// Assert
	require.Error(t, err)
	assert.Equal(t, ErrNoRPCPassword, err)
}

// TestLoadConfig_MissingRPCHost tests the failure case of LoadConfig when RPC_HOST is missing
func TestLoadConfig_MissingRPCHost(t *testing.T) {
	// Setup
	tester.SetupEnv(t)
	tester.UnsetEnv(t, "RPC_HOST")

	defer func() {
		tester.TeardownEnv(t)
	}()

	// Execute
	_, err := LoadConfig(context.Background(), nil, true)

	// Assert
	require.Error(t, err)
	assert.Equal(t, ErrNoRPCHost, err)
}

// TestLoadConfig_MissingP2PIP tests the failure case of LoadConfig when P2P_IP is missing
func TestLoadConfig_MissingP2PIP(t *testing.T) {
	// Setup
	tester.SetupEnv(t)
	tester.UnsetEnv(t, "P2P_IP")

	defer func() {
		tester.TeardownEnv(t)
	}()

	// Execute
	_, err := LoadConfig(context.Background(), nil, true)

	// Assert
	require.Error(t, err)
	assert.Equal(t, ErrNoP2PIP, err)
}

// TestLoadConfig_MissingP2PPort tests the failure case of LoadConfig when P2P_PORT is missing
func TestLoadConfig_MissingP2PPort(t *testing.T) {
	// Setup
	tester.SetupEnv(t)
	tester.UnsetEnv(t, "P2P_PORT")

	defer func() {
		tester.TeardownEnv(t)
	}()

	// Execute
	_, err := LoadConfig(context.Background(), nil, true)

	// Assert
	require.Error(t, err)
	assert.Equal(t, ErrNoP2PPort, err)
}

// TestLoadConfig_MissingP2PPrivateKeyPath tests the failure case of LoadConfig when P2P_PRIVATE_KEY_PATH is missing
func TestLoadConfig_MissingP2PPrivateKeyPath(t *testing.T) {
	// Setup
	tester.SetupEnv(t)
	tester.UnsetEnv(t, "P2P_PRIVATE_KEY_PATH")

	defer func() {
		tester.TeardownEnv(t)
	}()

	// Execute
	_, err := LoadConfig(context.Background(), nil, true)

	// Assert
	require.NoError(t, err)
}

// TestLoadConfig_OverrideP2PBootstrapPeer tests the case of LoadConfig when P2P_BOOTSTRAP_PEER is set
func TestLoadConfig_OverrideP2PBoostrapPeer(t *testing.T) {
	// Setup
	tester.SetupEnv(t)
	tester.SetEnv(t, "P2P_BOOTSTRAP_PEER", "foobar")

	defer func() {
		tester.TeardownEnv(t)
	}()

	// Execute
	c, err := LoadConfig(context.Background(), nil, true)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "foobar", c.P2PBootstrapPeer)
}

// TestLoadConfig_OverrideP2PAlertSystemProtocolID tests the case of LoadConfig when P2P_ALERT_SYSTEM_PROTOCOL_ID is set
func TestLoadConfig_OverrideP2PAlertSystemProtocolID(t *testing.T) {
	// Setup
	tester.SetupEnv(t)
	tester.SetEnv(t, "P2P_ALERT_SYSTEM_PROTOCOL_ID", "foobar/1.0.1")

	defer func() {
		tester.TeardownEnv(t)
	}()

	// Execute
	c, err := LoadConfig(context.Background(), nil, true)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "foobar/1.0.1", c.P2PAlertSystemProtocolID)
}

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
