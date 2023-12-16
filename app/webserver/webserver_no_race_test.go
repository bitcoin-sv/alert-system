//go:build !race

package webserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/stretchr/testify/require"
)

// TestServer_Shutdown will test the method Shutdown()
func TestServer_Shutdown_NoRace(t *testing.T) {
	// t.Parallel()

	// This method currently has an issue with a race condition on starting/stopping the webserver
	t.Run("app config and load server", func(t *testing.T) {

		// Set the ctx
		ctx := context.Background()

		// Set the env to test
		err := os.Setenv(config.EnvironmentKey, config.EnvironmentTest)
		require.NoError(t, err)

		// Load the config from env/json
		var dependencies *config.Config
		dependencies, err = config.LoadDependencies(ctx, models.BaseModels, true)
		require.NoError(t, err)
		require.NotNil(t, dependencies)

		// Sync a new server
		s := NewServer(dependencies)
		require.NotNil(t, s)

		// todo having an issue starting webserver and shutting down (in different routines)

		// Sync the webserver (in go routine for non-blocking)
		idleConnectionsClosed := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			// We received an interrupt signal, shut down.
			// ctx2, cancel := context.WithTimeout(context.Background(), config.DefaultServerShutdown)
			// defer cancel()

			// err2 := s.Shutdown(ctx2)
			// require.NoError(t, err2)

			close(idleConnectionsClosed)
		}()

		go func() {
			// Serve the server!
			s.Serve()
		}()

		// Delay
		t.Log("server loaded, waiting 1 second...")
		time.Sleep(1 * time.Second)

		// Kill process
		err = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		require.NoError(t, err)

		// Shutdown
		err = s.Shutdown(ctx)
		require.NoError(t, err)
	})
}
