package model

import (
	"bytes"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIDs_GormDataType will test the method GormDataType()
func TestIDs_GormDataType(t *testing.T) {
	t.Parallel()

	i := new(IDs)
	assert.Equal(t, gormTypeText, i.GormDataType())
}

// TestIDs_Scan will test the method Scan()
func TestIDs_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		i := IDs{}
		err := i.Scan(nil)
		require.NoError(t, err)
		assert.Empty(t, i)
	})

	t.Run("empty string", func(t *testing.T) {
		i := IDs{}
		err := i.Scan("\"\"")
		require.Error(t, err)
		assert.Empty(t, i)
	})

	t.Run("valid slice of ids", func(t *testing.T) {
		i := IDs{}
		err := i.Scan("[\"test1\",\"test2\"]")
		require.NoError(t, err)
		assert.Len(t, i, 2)
		assert.Equal(t, "test1", i[0])
		assert.Equal(t, "test2", i[1])
	})

	t.Run("empty id slice", func(t *testing.T) {
		i := IDs{}
		err := i.Scan("[\"\"]")
		require.NoError(t, err)
		assert.Len(t, i, 1)
		assert.Empty(t, i[0])
	})

	t.Run("invalid JSON", func(t *testing.T) {
		i := IDs{}
		err := i.Scan("[test1]")
		require.Error(t, err)
	})
}

// TestIDs_Value will test the method Value()
func TestIDs_Value(t *testing.T) {
	t.Parallel()

	t.Run("empty object", func(t *testing.T) {
		i := IDs{}
		value, err := i.Value()
		require.NoError(t, err)
		assert.Equal(t, "[]", value)
	})

	t.Run("ids present", func(t *testing.T) {
		i := IDs{"test1"}
		value, err := i.Value()
		require.NoError(t, err)
		assert.Equal(t, "[\"test1\"]", value)
	})
}

// TestMarshalIDs will test the method MarshalIDs()
func TestMarshalIDs(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		writer := MarshalIDs(nil)
		require.NotNil(t, writer)
		assert.IsType(t, graphql.Null, writer)
	})

	t.Run("empty object", func(t *testing.T) {
		writer := MarshalIDs(IDs{})
		require.NotNil(t, writer)
		b := bytes.NewBufferString("")
		writer.MarshalGQL(b)
		assert.Equal(t, "[]\n", b.String())
	})

	t.Run("map present", func(t *testing.T) {
		writer := MarshalIDs(IDs{"test1"})
		require.NotNil(t, writer)
		b := bytes.NewBufferString("")
		writer.MarshalGQL(b)
		assert.Equal(t, "[\"test1\"]\n", b.String())
	})
}

// TestUnmarshalIDs will test the method UnmarshalIDs()
func TestUnmarshalIDs(t *testing.T) {

	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		i, err := UnmarshalIDs(nil)
		require.NoError(t, err)
		assert.Empty(t, i)
		assert.IsType(t, IDs{}, i)
	})

	t.Run("empty string", func(t *testing.T) {
		i, err := UnmarshalIDs("\"\"")
		require.Error(t, err)
		assert.Empty(t, i)
		assert.IsType(t, IDs{}, i)
	})

	t.Run("valid set of ids", func(t *testing.T) {
		i, err := UnmarshalIDs(IDs{"test1", "test2"})
		require.NoError(t, err)
		assert.Len(t, i, 2)
		assert.IsType(t, IDs{}, i)
		assert.Equal(t, "test1", i[0])
		assert.Equal(t, "test2", i[1])
	})
}
