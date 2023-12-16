package webserver

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewServer will test the method NewServer()
func TestNewServer(t *testing.T) {
	t.Parallel()

	t.Run("empty values", func(t *testing.T) {
		s := NewServer(nil)
		require.NotNil(t, s)
		assert.Nil(t, s.Config)
		assert.Nil(t, s.Router)
		assert.Nil(t, s.WebServer)
	})

	t.Run("set values", func(t *testing.T) {
		dependencies := &config.Config{}
		s := NewServer(dependencies)
		require.NotNil(t, s)
		assert.Equal(t, dependencies, s.Config)
		assert.Equal(t, dependencies, s.Config)
		assert.Nil(t, s.Router)
		assert.Nil(t, s.WebServer)
	})
}

// TestServer_Shutdown will test the method Shutdown()
func TestServer_Shutdown(t *testing.T) {
	t.Parallel()

	t.Run("no server, services", func(t *testing.T) {
		s := NewServer(nil)
		require.NotNil(t, s)

		err := s.Shutdown(context.Background())
		require.NoError(t, err)
	})

	t.Run("basic app config and services", func(t *testing.T) {
		dependencies := &config.Config{}

		s := NewServer(dependencies)
		require.NotNil(t, s)

		err := s.Shutdown(context.Background())
		require.NoError(t, err)
	})

	t.Run("app config from json", func(t *testing.T) {

		// Set the ctx
		ctx := context.Background()

		tester.SetupEnv(t)
		defer func() {
			tester.TeardownEnv(t)
		}()

		// Execute
		appConfig, err := config.LoadConfig(ctx, nil, true)
		require.NoError(t, err)

		// Load the config from env/json
		require.NoError(t, err)
		require.NotNil(t, appConfig)

		// Sync a new server
		s := NewServer(appConfig)
		require.NotNil(t, s)

		// Shutdown the server
		err = s.Shutdown(ctx)
		require.NoError(t, err)
	})
}
