package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewNodeConfig creates a new NodeConfig struct
func TestNewNodeConfig(t *testing.T) {
	t.Run("test new node config", func(t *testing.T) {
		// Create a new node config
		node := NewNodeConfig("user", "pass", "host")

		val := node.GetRPCUser()
		assert.Equal(t, "user", val)

		val = node.GetRPCPassword()
		assert.Equal(t, "pass", val)

		val = node.GetRPCHost()
		assert.Equal(t, "host", val)
	})
}

// TestNewNodeMock creates a new NodeConfig struct for testing
func TestNewNodeMock(t *testing.T) {
	t.Run("test new node mock", func(t *testing.T) {
		// Create a new node config
		node := NewNodeMock("user", "pass", "host")

		val := node.GetRPCUser()
		assert.Equal(t, "user", val)

		val = node.GetRPCPassword()
		assert.Equal(t, "pass", val)

		val = node.GetRPCHost()
		assert.Equal(t, "host", val)
	})
}
