// Package app provides the overarching application logic
package app

import (
	"net/http"

	"github.com/bitcoin-sv/alert-system/app/config"
	apirouter "github.com/mrz1836/go-api-router"
)

// Action is the configuration for the actions and related services
type Action struct {
	Config *config.Config // Combination of configuration and services, being passed down into the handlers
}

// APIError is the enriched error message for API related errors
type APIError struct {
	Message    string `json:"message" url:"message"`         // Public error message
	StatusCode int    `json:"status_code" url:"status_code"` // Associated HTTP status code (should be in request as well)
}

// APIErrorResponse will return an error response message
func APIErrorResponse(w http.ResponseWriter, req *http.Request, statusCode int, err error) {
	apirouter.ReturnResponse(
		w, req, statusCode,
		&APIError{
			Message:    err.Error(),
			StatusCode: statusCode,
		},
	)
}
