package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestName_IsEmpty will test the method IsEmpty()
func TestName_IsEmpty(t *testing.T) {
	t.Parallel()

	t.Run("valid empty name", func(t *testing.T) {
		n := NameEmpty
		require.True(t, n.IsEmpty())
	})

	t.Run("empty string", func(t *testing.T) {
		n := new(Name)
		require.True(t, n.IsEmpty())
	})

	t.Run("valid name", func(t *testing.T) {
		n := NameAlertMessage
		require.False(t, n.IsEmpty())
	})
}
