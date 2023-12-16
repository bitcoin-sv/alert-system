package base

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIndex will test the method index()
func (ts *TestSuite) TestIndex() {
	ts.T().Run("test index", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// Fire the request
		index(w, req, nil)
		res := w.Result()
		defer func() {
			_ = res.Body.Close()
		}()

		// Test the body
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "\"Bitcoin SV Alert System\"\n", string(data))

		// Check the result
		require.Equal(t, "200 OK", res.Status)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}
