// Package tester is for testing the alert system
package tester

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// SetEnv will set an environment variable
func SetEnv(t *testing.T, key, value string) {
	err := os.Setenv(key, value)
	require.NoError(t, err)
}

// UnsetEnv will unset an environment variable
func UnsetEnv(t *testing.T, key string) {
	err := os.Unsetenv(key)
	require.NoError(t, err)
}

// SetupEnv helper function to set up environment for testing
func SetupEnv(t *testing.T) {
	SetEnv(t, "RPC_USER", "user")
	SetEnv(t, "RPC_PASSWORD", "password")
	SetEnv(t, "RPC_HOST", "localhost")
	SetEnv(t, "P2P_PRIVATE_KEY_PATH", "/path/to/private/key")
	SetEnv(t, "P2P_IP", "192.168.1.1")
	SetEnv(t, "P2P_PORT", "8000")
	SetEnv(t, "ALERT_WEBHOOK_URL", "https://webhook.url")
	SetEnv(t, "DATABASE_PATH", "")
}

// TeardownEnv helper function to tear down environment after testing
func TeardownEnv(t *testing.T) {
	UnsetEnv(t, "RPC_USER")
	UnsetEnv(t, "RPC_PASSWORD")
	UnsetEnv(t, "RPC_HOST")
	UnsetEnv(t, "P2P_PRIVATE_KEY_PATH")
	UnsetEnv(t, "P2P_IP")
	UnsetEnv(t, "P2P_PORT")
	UnsetEnv(t, "ALERT_WEBHOOK_URL")
	UnsetEnv(t, "DATABASE_PATH")
}
