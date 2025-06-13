package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModelSetRecordTime will test the method SetRecordTime()
func TestModelSetRecordTime(t *testing.T) {
	t.Parallel()

	t.Run("empty model", func(t *testing.T) {
		m := new(Model)
		assert.True(t, m.CreatedAt.IsZero())
		assert.True(t, m.UpdatedAt.IsZero())
	})

	t.Run("set created at time", func(t *testing.T) {
		m := new(Model)
		m.SetRecordTime(true)
		assert.False(t, m.CreatedAt.IsZero())
		assert.False(t, m.UpdatedAt.IsZero())
		assert.Equal(t, m.CreatedAt.Hour(), m.UpdatedAt.Hour())   // Both should be the same time when created
		assert.Equal(t, m.CreatedAt.Month(), m.UpdatedAt.Month()) // Both should be the same time when created
		assert.Equal(t, m.CreatedAt.Year(), m.UpdatedAt.Year())   // Both should be the same time when created
		assert.Equal(t, m.CreatedAt.Day(), m.UpdatedAt.Day())     // Both should be the same time when created
	})

	t.Run("set updated at time", func(t *testing.T) {
		m := new(Model)
		m.SetRecordTime(false)
		assert.True(t, m.CreatedAt.IsZero())
		assert.False(t, m.UpdatedAt.IsZero())
	})

	t.Run("set both times", func(t *testing.T) {
		m := new(Model)
		m.SetRecordTime(false)
		m.SetRecordTime(true)
		assert.False(t, m.CreatedAt.IsZero())
		assert.False(t, m.UpdatedAt.IsZero())
	})
}

// TestModelSetRecordTime will test the method New()
func TestModelNew(t *testing.T) {
	t.Parallel()

	t.Run("New model", func(t *testing.T) {
		m := new(Model)
		assert.False(t, m.IsNew())
	})

	t.Run("set New flag", func(t *testing.T) {
		m := new(Model)
		m.New()
		assert.True(t, m.IsNew())
	})
}

// TestModelGetOptions will test the method GetOptions()
func TestModelGetOptions(t *testing.T) {
	// t.Parallel()

	t.Run("base model", func(t *testing.T) {
		m := new(Model)
		opts := m.GetOptions(false)
		assert.Empty(t, opts)
	})

	t.Run("new record model", func(t *testing.T) {
		m := new(Model)
		opts := m.GetOptions(true)
		assert.Len(t, opts, 1)
	})
}

// TestModel_IsNew will test the method IsNew()
func TestModel_IsNew(t *testing.T) {
	t.Parallel()

	t.Run("base model", func(t *testing.T) {
		m := new(Model)
		assert.False(t, m.IsNew())
	})

	t.Run("New model", func(t *testing.T) {
		m := new(Model)
		m.New()
		assert.True(t, m.IsNew())
	})
}

// TestModel_Name will test the method Name()
func TestModel_Name(t *testing.T) {
	t.Parallel()

	t.Run("base model", func(t *testing.T) {
		m := new(Model)
		assert.Empty(t, m.Name())
	})

	t.Run("set model name - alert message", func(t *testing.T) {
		m := new(Model)
		m.name = NameAlertMessage
		assert.Equal(t, "alert_message", m.Name())
	})

	t.Run("set model name - public key", func(t *testing.T) {
		m := new(Model)
		m.name = NamePublicKey
		assert.Equal(t, "public_key", m.Name())
	})
}

// TestModel_GetID will test the method GetID()
func TestModel_GetID(t *testing.T) {
	t.Parallel()

	t.Run("base model - no id", func(t *testing.T) {
		m := new(Model)
		assert.Empty(t, m.GetID())
	})
}

// TestModel_GetTableName will test the method GetTableName()
func TestModel_GetTableName(t *testing.T) {
	t.Parallel()

	t.Run("base model", func(t *testing.T) {
		m := new(Model)
		assert.Equal(t, TableEmpty, m.GetTableName())
	})
}

// TestModel_NotNew will test the method NotNew()
func TestModel_NotNew(t *testing.T) {
	t.Parallel()

	t.Run("base model", func(t *testing.T) {
		m := new(Model)
		assert.False(t, m.IsNew())
	})

	t.Run("set: new and not new", func(t *testing.T) {
		m := new(Model)
		m.SetOptions(New())
		assert.True(t, m.IsNew())

		m.NotNew()
		assert.False(t, m.IsNew())
	})
}

// TestModel_AfterCreated will test the method AfterCreated()
func TestModel_AfterCreated(t *testing.T) {
	t.Parallel()

	t.Run("no customization, no debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		err := m.AfterCreated(context.Background())
		require.NoError(t, err)
	})

	t.Run("no customization, with debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		m.Debug(true)
		err := m.AfterCreated(context.Background())
		require.NoError(t, err)
	})
}

// TestModel_AfterDeleted will test the method AfterDeleted()
func TestModel_AfterDeleted(t *testing.T) {
	t.Parallel()

	t.Run("no customization, no debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		err := m.AfterDeleted(context.Background())
		require.NoError(t, err)
	})

	t.Run("no customization, with debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		m.Debug(true)
		err := m.AfterDeleted(context.Background())
		require.NoError(t, err)
	})
}

// TestModel_AfterUpdated will test the method AfterUpdated()
func TestModel_AfterUpdated(t *testing.T) {
	t.Parallel()

	t.Run("no customization, no debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		err := m.AfterUpdated(context.Background())
		require.NoError(t, err)
	})

	t.Run("no customization, with debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		m.Debug(true)
		err := m.AfterUpdated(context.Background())
		require.NoError(t, err)
	})
}

// TestModel_BeforeCreating will test the method BeforeCreating()
func TestModel_BeforeCreating(t *testing.T) {
	t.Parallel()

	t.Run("no customization, no debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		err := m.BeforeCreating(context.Background())
		require.NoError(t, err)
	})

	t.Run("no customization, with debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		m.Debug(true)
		err := m.BeforeCreating(context.Background())
		require.NoError(t, err)
	})
}

// TestModel_BeforeUpdating will test the method BeforeUpdating()
func TestModel_BeforeUpdating(t *testing.T) {
	t.Parallel()

	t.Run("no customization, no debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		err := m.BeforeUpdating(context.Background())
		require.NoError(t, err)
	})

	t.Run("no customization, with debug", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		m.Debug(true)
		err := m.BeforeUpdating(context.Background())
		require.NoError(t, err)
	})
}

// TestModel_ChildModels will test the method ChildModels()
func TestModel_ChildModels(t *testing.T) {
	t.Parallel()

	m := NewBaseModel(NameEmpty)
	require.NotNil(t, m)
	models := m.ChildModels()
	require.Nil(t, models)
}

// TestModel_Debug will test the method Debug()
func TestModel_Debug(t *testing.T) {
	t.Parallel()

	t.Run("enabled", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		require.NotNil(t, m)
		m.Debug(true)
		assert.True(t, m.debug)
	})

	t.Run("disabled", func(t *testing.T) {
		m := NewBaseModel(NameEmpty)
		require.NotNil(t, m)
		m.Debug(false)
		assert.False(t, m.debug)
	})
}

// TestModel_Enrich will test the method enrich()
func TestModel_Enrich(t *testing.T) {
	t.Parallel()

	n := Model{
		name: NameEmpty,
	}

	user := Model{
		name:  NameAlertMessage,
		debug: true,
	}

	n.enrich(user.name, user.GetOptions(true)...)
	require.True(t, n.debug)
	require.True(t, n.IsNew())
	require.Equal(t, NameAlertMessage.String(), n.Name())
}
