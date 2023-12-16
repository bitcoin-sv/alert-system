package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestHead will test the method Head()
func TestHead(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Fire the request
	Head(w, req, nil)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	// Check body
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Empty(t, data)

	// Check the result
	require.Equal(t, "200 OK", res.Status)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

// TestNotFound will test the method NotFound()
func TestNotFound(t *testing.T) {
	path := "/unknown-target"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()

	// Fire the request
	NotFound(w, req)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	// Check body
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "\"route not found: "+path+"\"\n", string(data))

	// Check the result
	require.Equal(t, "404 Not Found", res.Status)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

// TestMethodNotAllowed will test the method MethodNotAllowed()
func TestMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodConnect, "/", nil)
	w := httptest.NewRecorder()

	// Fire the request
	MethodNotAllowed(w, req)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	// Check body
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "\"/:"+http.MethodConnect+"\"\n", string(data))

	// Check the result
	require.Equal(t, "405 Method Not Allowed", res.Status)
	require.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}
