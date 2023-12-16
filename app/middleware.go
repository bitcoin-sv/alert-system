package app

import (
	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// NewStack is used for registering routes
func NewStack(conf *config.Config) (Action, *apirouter.InternalStack) {
	return Action{Config: conf}, apirouter.NewStack()
}

// Request will process the request in the router
// This is used for logging requests or not logging the requests from the router
func (a *Action) Request(router *apirouter.Router, h httprouter.Handle) httprouter.Handle {
	if a.Config.RequestLogging {
		return router.Request(h)
	}
	return router.RequestNoLogging(h)
}
