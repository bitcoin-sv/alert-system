package tester

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSetupEnv will test the method SetupEnv()
func TestSetEnv(t *testing.T) {
	// Define a key-value pair for the test
	key := "TEST_ENV_VAR"
	value := "test_value"

	// Call SetEnv to set the environment variable
	SetEnv(t, key, value)

	// Retrieve the value of the environment variable
	setValue, exists := os.LookupEnv(key)
	require.True(t, exists, "Environment variable should exist")
	require.Equal(t, value, setValue, "Environment variable should have the correct value")

	// Clean up
	err := os.Unsetenv(key)
	require.NoError(t, err, "Unsetting environment variable should not produce an error")
}

// TestUnsetEnv will test the method UnsetEnv()
func TestUnsetEnv(t *testing.T) {
	// Define a key for the test and set the environment variable
	key := "TEST_ENV_VAR"
	value := "test_value"
	err := os.Setenv(key, value)
	require.NoError(t, err, "Setting environment variable should not produce an error")

	// Call UnsetEnv to unset the environment variable
	UnsetEnv(t, key)

	// Check if the environment variable still exists
	_, exists := os.LookupEnv(key)
	require.False(t, exists, "Environment variable should not exist after unsetting")
}

// TestSetupEnv will test the method SetupEnv()
func TestSetupEnv(t *testing.T) {
	// Call SetupEnv to set up the environment variables
	SetupEnv(t)

	defer func() {
		TeardownEnv(t)
	}()

	// Define a map of expected environment variables and their values
	expectedVars := map[string]string{
		"RPC_USER":             "user",
		"RPC_PASSWORD":         "password",
		"RPC_HOST":             "localhost",
		"P2P_PRIVATE_KEY_PATH": "/path/to/private/key",
		"P2P_IP":               "192.168.1.1",
		"P2P_PORT":             "8000",
		"ALERT_WEBHOOK_URL":    "https://webhook.url",
		"DATABASE_PATH":        "",
	}

	// Check each environment variable
	for key, expectedValue := range expectedVars {
		value, exists := os.LookupEnv(key)
		require.True(t, exists, "Environment variable %s should exist", key)
		require.Equal(t, expectedValue, value, "Environment variable %s should have the correct value", key)
	}
}

// TestTeardownEnv will test the method TeardownEnv()
func TestTeardownEnv(t *testing.T) {
	// First, call SetupEnv to set up the environment variables
	SetupEnv(t)

	// Call TeardownEnv to remove the environment variables
	TeardownEnv(t)

	// List all environment variables that should have been removed
	envVars := []string{
		"RPC_USER",
		"RPC_PASSWORD",
		"RPC_HOST",
		"P2P_PRIVATE_KEY_PATH",
		"P2P_IP",
		"P2P_PORT",
		"ALERT_WEBHOOK_URL",
		"DATABASE_PATH",
	}

	// Check each environment variable to ensure it has been removed
	for _, key := range envVars {
		_, exists := os.LookupEnv(key)
		require.False(t, exists, "Environment variable %s should not exist after teardown", key)
	}
}
