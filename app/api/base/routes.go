package base

import (
	"net/http"

	"github.com/bitcoin-sv/alert-system/app/p2p"

	"github.com/bitcoin-sv/alert-system/app"
	"github.com/bitcoin-sv/alert-system/app/config"
	apirouter "github.com/mrz1836/go-api-router"
)

// Action is an extension of app.Action for this package
type Action struct {
	app.Action
}

// RegisterRoutes register all the package specific routes
func RegisterRoutes(router *apirouter.Router, conf *config.Config, p2pServ *p2p.Server) {

	// Load the actions and set the services
	action := &Action{app.Action{Config: conf, P2pServer: p2pServ}}

	// Set the main index page (navigating to slash or the root of the major version)
	router.HTTPRouter.GET("/", action.Request(router, action.index))

	// Options request (for CORs)
	router.HTTPRouter.OPTIONS("/", router.SetCrossOriginHeaders)

	// Head requests are sometimes used for CORs
	router.HTTPRouter.HEAD("/", app.Head)

	// Set the 404 handler (any request not detected)
	router.HTTPRouter.NotFound = http.HandlerFunc(app.NotFound)

	// Set the method not allowed
	router.HTTPRouter.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowed)

	// Set the health request
	router.HTTPRouter.GET("/health", action.Request(router, action.health))

	// Set the get alerts request
	router.HTTPRouter.GET("/alerts", action.Request(router, action.alerts))

	// Set the get alert request
	router.HTTPRouter.GET("/alert/:sequence", action.Request(router, action.alert))
}
