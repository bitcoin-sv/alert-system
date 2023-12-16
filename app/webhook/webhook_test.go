package webhook

import (
	"errors"
	"net/http"
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
/*func TestPostAlert(t *testing.T) {
	// Create a mock HTTP server for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a successful response from the webhook server
		w.WriteHeader(http.StatusOK)
	}))

	defer mockServer.Close()

	// Create a mock alert message for testing
	mockAlert := &models.AlertMessage{
		// Set your alert message fields here
		Raw: "01000000150000005247bd6500000000010000000e546869732069732061207465737420bd1521c60845302ca088f8626ce77cef64e65b21f09de1cd2aa466e774421d61310141628fa14478af8c8134540b08149db916085f8d61c0277b8b9f1473c0161fb79c0667e48af7fefcdb963673c5a03546f7885ece9b4d2fb44138eee3c53ed055a575872fc3f93afad934abd77038d5f546df639259e9b5192bdcedc036f6b61f51312c120d76e5031709a9b03dc52ef4e8198eb4591703d5c2a56cc2c1960e5c1aeb792acbd68d3c0bd2f3000345a0d6b979a276068ef24ffafd33c22eba01ef",
	}
	mockAlert.SetAlertType(models.AlertTypeInformational)
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
}*/
