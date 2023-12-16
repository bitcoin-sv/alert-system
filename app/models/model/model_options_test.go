package model

import (
	"testing"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/ordishs/gocore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew will test the method New()
func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("Get opts", func(t *testing.T) {
		opt := New()
		assert.IsType(t, *new(Options), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opt := New()
		m := new(Model)
		m.SetOptions(opt)
		assert.True(t, m.IsNew())
	})
}

// TestWithMetadata will test the method WithMetadata()
func TestWithMetadata(t *testing.T) {
	t.Parallel()

	t.Run("Get opts", func(t *testing.T) {
		opt := WithMetadata("key", "value")
		assert.IsType(t, *new(Options), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opt := WithMetadata("key", "value")
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, "value", m.Metadata.GetKey("key"))
	})
}

// TestWithMetadatas will test the method WithMetadatas()
func TestWithMetadatas(t *testing.T) {
	t.Parallel()

	t.Run("Get opts", func(t *testing.T) {
		opt := WithMetadatas(map[string]interface{}{
			"key": "value",
		})
		assert.IsType(t, *new(Options), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opt := WithMetadatas(map[string]interface{}{
			"key": "value",
		})
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, "value", m.Metadata.GetKey("key"))
	})
}

// TestWithLogger will test the method WithLogger()
func TestWithLogger(t *testing.T) {
	t.Parallel()

	t.Run("general opts", func(t *testing.T) {
		opt := WithLogger(nil)
		assert.IsType(t, *new(Options), opt)
	})

	t.Run("nil logger", func(t *testing.T) {
		opt := WithLogger(nil)
		m := new(Model)
		m.SetOptions(opt)
		assert.Nil(t, m.Logger())
	})

	t.Run("valid logger", func(t *testing.T) {
		l := &config.ExtendedLogger{
			Logger: gocore.Log("test"),
		}
		require.NotNil(t, l)
		opt := WithLogger(l)
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, l, m.Logger())
	})
}

// TestWithDebug will test the method WithDebug()
func TestWithDebug(t *testing.T) {
	t.Parallel()

	t.Run("general opts", func(t *testing.T) {
		opt := WithDebug()
		assert.IsType(t, *new(Options), opt)
	})

	t.Run("enable debug", func(t *testing.T) {
		opt := WithDebug()
		m := new(Model)
		m.SetOptions(opt)
		assert.True(t, m.debug)
	})
}

// TestWithAllDependencies will test the method WithAllDependencies()
func TestWithAllDependencies(t *testing.T) {
	t.Parallel()

	t.Run("general opts", func(t *testing.T) {
		opt := WithAllDependencies(nil)
		assert.IsType(t, *new(Options), opt)
	})

	t.Run("nil dependencies", func(t *testing.T) {
		opt := WithAllDependencies(nil)
		m := new(Model)
		m.SetOptions(opt)

		assert.Panics(t, func() {
			assert.Nil(t, m.Datastore())
		})
	})

	/*t.Run("valid dependencies", func(t *testing.T) {
		err := os.Setenv(config.EnvironmentKey, config.EnvironmentTest)
		require.NoError(t, err)

		var ad *config.AppDependencies
		ad, err = config.LoadDependencies(context.Background(), "../../", nil)
		require.NoError(t, err)
		require.NotNil(t, ad)

		opt := WithAllDependencies(ad)
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, ad.AppConfig, m.Config())
		assert.Equal(t, ad.AppServices.Cachestore, m.Cachestore())
		assert.Equal(t, ad.AppServices.Datastore, m.Datastore())
	})*/
}
