package config

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mrz1836/go-datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig_Failure tests the success case of LoadDependencies
func TestLoadConfig_Success(t *testing.T) {
	t.Run("successfully loading the config", func(t *testing.T) {
		err := os.Setenv(EnvironmentKey, EnvironmentTest)
		require.NoError(t, err)

		// Execute
		var c *Config
		c, err = LoadDependencies(context.Background(), nil, true)
		require.NoError(t, err)

		defer c.CloseAll(context.Background())

		// Assert
		assert.NotNil(t, c)
		assert.Equal(t, "/path/to/private/key", c.P2P.PrivateKeyPath)
		assert.Empty(t, c.P2P.BootstrapPeer)
		assert.Equal(t, DefaultAlertSystemProtocolID, c.P2P.AlertSystemProtocolID)
		assert.Equal(t, DefaultPeerDiscoveryInterval, c.P2P.PeerDiscoveryInterval)
		assert.Equal(t, DefaultAlertProcessingInterval, c.AlertProcessingInterval)
		assert.Equal(t, "192.168.1.1", c.P2P.IP)
		assert.Equal(t, "8000", c.P2P.Port)
		assert.Equal(t, "https://webhook.url", c.AlertWebhookURL)
	})
}

// TestLoadConfigFile tests the method LoadConfigFile()
func TestLoadConfigFile(t *testing.T) {

	t.Run("no env", func(t *testing.T) {
		err := os.Unsetenv(EnvironmentKey)
		require.NoError(t, err)

		var ac *Config
		ac, err = LoadConfigFile()
		require.Error(t, err)
		require.Nil(t, ac)
		assert.Contains(t, err.Error(), "invalid environment")
	})

	t.Run("missing rpc connections", func(t *testing.T) {
		err := os.Setenv(EnvironmentKey, EnvironmentTest)
		require.NoError(t, err)

		// err = os.Setenv("ALERT_SYSTEM_RPC_CONNECTIONS", "[{\"user\":\"galt\",\"password\":\"galt\",\"host\":\"http://localhost:8333\"}]")
		err = os.Setenv("ALERT_SYSTEM_RPC_CONNECTIONS", "[]")
		require.NoError(t, err)
		defer func() {
			_ = os.Unsetenv("ALERT_SYSTEM_RPC_CONNECTIONS")
		}()

		// Execute
		var c *Config
		c, err = LoadDependencies(context.Background(), nil, true)
		require.Nil(t, c)
		require.Error(t, err)
	})

	t.Run("missing ip address", func(t *testing.T) {
		err := os.Setenv(EnvironmentKey, EnvironmentTest)
		require.NoError(t, err)

		err = os.Setenv("ALERT_SYSTEM_P2P__IP", " ")
		require.NoError(t, err)
		defer func() {
			_ = os.Unsetenv("ALERT_SYSTEM_P2P__IP")
		}()

		// Execute
		var c *Config
		c, err = LoadDependencies(context.Background(), nil, true)
		require.Nil(t, c)

		require.Error(t, err)
		assert.Equal(t, ErrNoP2PIP, err)
	})

	t.Run("missing port", func(t *testing.T) {
		err := os.Setenv(EnvironmentKey, EnvironmentTest)
		require.NoError(t, err)

		err = os.Setenv("ALERT_SYSTEM_P2P__PORT", " ")
		require.NoError(t, err)
		defer func() {
			_ = os.Unsetenv("ALERT_SYSTEM_P2P__PORT")
		}()

		// Execute
		var c *Config
		c, err = LoadDependencies(context.Background(), nil, true)
		require.Nil(t, c)

		require.Error(t, err)
		assert.Equal(t, ErrNoP2PPort, err)
	})

	t.Run("invalid custom file path for config", func(t *testing.T) {
		err := os.Setenv(EnvironmentKey, EnvironmentTest)
		require.NoError(t, err)

		err = os.Setenv(EnvironmentCustomFilePath, "file-not-found.json")
		require.NoError(t, err)
		defer func() {
			_ = os.Unsetenv(EnvironmentCustomFilePath)
		}()

		// Execute
		var c *Config
		c, err = LoadDependencies(context.Background(), nil, true)
		require.Nil(t, c)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("valid custom location for config file", func(t *testing.T) {
		err := os.Setenv(EnvironmentKey, EnvironmentTest)
		require.NoError(t, err)

		err = os.Setenv(EnvironmentCustomFilePath, "envs/test.json")
		require.NoError(t, err)
		defer func() {
			_ = os.Unsetenv(EnvironmentCustomFilePath)
		}()

		// Execute
		var c *Config
		c, err = LoadDependencies(context.Background(), nil, true)
		require.NotNil(t, c)
		require.NoError(t, err)
		defer c.CloseAll(context.Background())
	})

	t.Run("test env, found file, test all structs", func(t *testing.T) {
		err := os.Setenv(EnvironmentKey, EnvironmentTest)
		require.NoError(t, err)

		var ac *Config
		ac, err = LoadConfigFile()
		require.NoError(t, err)
		require.NotNil(t, ac)

		defer ac.CloseAll(context.Background())

		assert.True(t, ac.RequestLogging)

		// Check nested structs (Webserver)
		assert.Equal(t, 60*time.Second, ac.WebServer.IdleTimeout)
		assert.Equal(t, 15*time.Second, ac.WebServer.ReadTimeout)
		assert.Equal(t, 15*time.Second, ac.WebServer.WriteTimeout)
		assert.Equal(t, "3000", ac.WebServer.Port)

		// Check nested structs (Datastore)
		assert.True(t, ac.Datastore.AutoMigrate)
		assert.True(t, ac.Datastore.Debug)
		assert.Equal(t, datastore.SQLite, ac.Datastore.Engine)
		assert.Empty(t, ac.Datastore.Password)
		assert.Equal(t, "alert_system", ac.Datastore.TablePrefix)
		assert.Empty(t, ac.Datastore.SQLite.DatabasePath)
		assert.False(t, ac.Datastore.SQLite.Shared)
		assert.Equal(t, "postgresql", ac.Datastore.SQLRead.Driver)
		assert.Equal(t, "localhost", ac.Datastore.SQLRead.Host)
		assert.Equal(t, time.Duration(20000000000), ac.Datastore.SQLRead.MaxConnectionIdleTime)
		assert.Equal(t, time.Duration(20000000000), ac.Datastore.SQLRead.MaxConnectionTime)
		assert.Equal(t, 2, ac.Datastore.SQLRead.MaxIdleConnections)
		assert.Equal(t, 5, ac.Datastore.SQLRead.MaxOpenConnections)
		assert.Equal(t, "postgresql", ac.Datastore.SQLWrite.Driver)
		assert.Equal(t, "localhost", ac.Datastore.SQLWrite.Host)
		assert.Equal(t, time.Duration(20000000000), ac.Datastore.SQLWrite.MaxConnectionIdleTime)
		assert.Equal(t, time.Duration(20000000000), ac.Datastore.SQLWrite.MaxConnectionTime)
		assert.Equal(t, 2, ac.Datastore.SQLWrite.MaxIdleConnections)
		assert.Equal(t, 5, ac.Datastore.SQLWrite.MaxOpenConnections)

		// RPC Connections
		assert.Len(t, ac.RPCConnections, 1)
		assert.Equal(t, "galt", ac.RPCConnections[0].User)
		assert.Equal(t, "galt", ac.RPCConnections[0].Password)
		assert.Equal(t, "http://localhost:8333", ac.RPCConnections[0].Host)
	})
}

// TestIsValidEnvironment will test the method isValidEnvironment()
func TestIsValidEnvironment(t *testing.T) {
	t.Run("empty env", func(t *testing.T) {
		valid := isValidEnvironment("")
		assert.False(t, valid)
	})

	t.Run("unknown env", func(t *testing.T) {
		valid := isValidEnvironment("unknown")
		assert.False(t, valid)
	})

	t.Run("different case of letters", func(t *testing.T) {
		valid := isValidEnvironment("LOCal")
		assert.True(t, valid)
	})

	t.Run("valid envs", func(t *testing.T) {
		valid := isValidEnvironment(EnvironmentTest)
		assert.True(t, valid)

		valid = isValidEnvironment(EnvironmentLocal)
		assert.True(t, valid)

		valid = isValidEnvironment(EnvironmentProduction)
		assert.True(t, valid)
	})
}
