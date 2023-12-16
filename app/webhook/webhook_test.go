package webhook

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockHTTPClient is a mock HTTP client for testing purposes
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do is the mock HTTP client Do function
func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.DoFunc != nil {
		return c.DoFunc(req)
	}
	return nil, errors.New("unimplemented")
}

// TestPostAlert tests the PostAlert function
func TestPostAlert(t *testing.T) {
	// Create a mock HTTP server for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a successful response from the webhook server
		w.WriteHeader(http.StatusOK)
	}))

	defer mockServer.Close()

	// Create a mock alert message for testing
	mockAlert := &models.AlertMessage{
		// Set your alert message fields here
	}

	t.Run("ValidPostAlert", func(t *testing.T) {
		// Initialize the PostAlert function with a mock HTTP client
		httpClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Ensure that the request is made to the mock server URL
				assert.Equal(t, mockServer.URL, req.URL.String())
				// Add more assertions as needed
				return nil, nil // Simulate a successful request
			},
		}

		// Call the PostAlert function with the mock HTTP client
		err := PostAlert(context.Background(), httpClient, mockServer.URL, mockAlert)

		// Check if there are no errors
		require.NoError(t, err)
	})

	t.Run("InvalidURL", func(t *testing.T) {
		// Call the PostAlert function with an empty URL
		err := PostAlert(context.Background(), nil, "", mockAlert)

		// Check if it returns an error for an empty URL
		require.Error(t, err)
		// Add more assertions as needed
	})

	t.Run("InvalidURLPrefix", func(t *testing.T) {
		// Call the PostAlert function with an invalid URL prefix
		err := PostAlert(context.Background(), nil, "invalid-url", mockAlert)

		// Check if it returns an error for an invalid URL prefix
		require.Error(t, err)
		// Add more assertions as needed
	})

	t.Run("EmptyURL", func(t *testing.T) {
		err := PostAlert(context.Background(), nil, "", mockAlert)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "webhook URL is not configured")
	})

	t.Run("InvalidURLPrefix", func(t *testing.T) {
		err := PostAlert(context.Background(), nil, "invalid-url", mockAlert)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "webhook URL [invalid-url] is does not have a valid prefix")
	})

	t.Run("HTTPClientError", func(t *testing.T) {
		httpClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("HTTP client error")
			},
		}

		err := PostAlert(context.Background(), httpClient, mockServer.URL, mockAlert)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP client error")
	})

	t.Run("InvalidResponseStatus", func(t *testing.T) {
		httpClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest, // Simulate a non-OK response status
				}, nil
			},
		}

		err := PostAlert(context.Background(), httpClient, mockServer.URL, mockAlert)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code [400] sending payload to webhook")
	})
}
