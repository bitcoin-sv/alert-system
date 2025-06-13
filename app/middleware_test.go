package app

import (
	"net/http"
	"testing"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/stretchr/testify/require"
)

func testHandle(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	apirouter.RespondWith(w, req, http.StatusOK, "")
}

// TestNewStack will test the method NewStack()
func TestNewStack(t *testing.T) {
	t.Parallel()

	dep := new(config.Config)

	a, stack := NewStack(dep)
	require.NotNil(t, a)
	require.NotNil(t, stack)
	require.Equal(t, dep, a.Config)
	require.IsType(t, &apirouter.InternalStack{}, stack)
}

// TestAction_Request will test the method Request()
func TestAction_Request(t *testing.T) {
	t.Parallel()

	t.Run("request logging: enabled", func(t *testing.T) {
		dep := new(config.Config)
		dep.RequestLogging = true
		a, stack := NewStack(dep)
		require.NotNil(t, a)
		require.NotNil(t, stack)

		router := apirouter.New()
		a.Request(router, testHandle)
	})

	t.Run("request logging: disabled", func(t *testing.T) {
		dep := new(config.Config)
		dep.RequestLogging = false
		a, stack := NewStack(dep)
		require.NotNil(t, a)
		require.NotNil(t, stack)

		router := apirouter.New()
		a.Request(router, testHandle)
	})
}
