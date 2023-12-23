package config

import (
	"context"
	"testing"

	"github.com/mrz1836/go-datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadDatastore tests the cases of loadDatastore
func TestLoadDatastore(t *testing.T) {

	t.Run("failure - datastore unsupported", func(t *testing.T) {
		// Execute
		c := &Config{
			Datastore: DatastoreConfig{
				Engine: "unsupported",
			},
		}
		err := c.loadDatastore(context.Background(), nil)

		// Assert
		require.Error(t, err)
		assert.Equal(t, ErrDatastoreUnsupported, err)
	})

	t.Run("success - sqlite", func(t *testing.T) {

		// Execute
		c := &Config{
			Services: Services{},
			Datastore: DatastoreConfig{
				Engine:      datastore.SQLite,
				AutoMigrate: true,
				TablePrefix: "test",
				Debug:       false,
				SQLite: &datastore.SQLiteConfig{
					CommonConfig: datastore.CommonConfig{
						Debug:              true,
						MaxIdleConnections: 1,
						MaxOpenConnections: 1,
					},
					Shared:       false,
					DatabasePath: "",
				},
			},
		}
		err := c.loadDatastore(context.Background(), nil)

		// Assert
		require.NoError(t, err)
	})
}
