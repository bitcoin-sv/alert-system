// Package webserver is the web HTTP server for the alert-system
package webserver

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"strings"

	"github.com/bitcoin-sv/alert-system/app/p2p"

	"github.com/bitcoin-sv/alert-system/app/api/base"
	"github.com/bitcoin-sv/alert-system/app/config"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
)

const (
	wildcard = "*"
)

// Server is the configuration, services, and actual web server
type Server struct {
	Config    *config.Config
	Router    *apirouter.Router
	WebServer *http.Server
	P2pServer *p2p.Server
}

// NewServer will return a new server service
func NewServer(conf *config.Config, serv *p2p.Server) *Server {
	return &Server{
		Config:    conf,
		P2pServer: serv,
	}
}

// Serve will load a server and start serving
func (s *Server) Serve() {

	// Load the server defaults
	s.WebServer = &http.Server{
		Addr:              ":" + s.Config.WebServer.Port,
		Handler:           s.Handlers(),
		IdleTimeout:       s.Config.WebServer.IdleTimeout,
		ReadHeaderTimeout: s.Config.WebServer.ReadTimeout,
		ReadTimeout:       s.Config.WebServer.ReadTimeout,
		WriteTimeout:      s.Config.WebServer.WriteTimeout,
		TLSConfig: &tls.Config{
			NextProtos:       []string{"h2", "http/1.1"},
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	// Turn off keep alive
	// s.WebServer.SetKeepAlivesEnabled(false)

	// Listen and serve
	if err := s.WebServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Config.Services.Log.Info("shutting down web server [" + err.Error() + "]...")
	}
}

// Shutdown will stop the web server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.Config != nil {
		s.Config.CloseAll(ctx) // Should have been executed in main.go, but might panic and not run?
	}
	if s.WebServer != nil {
		return s.WebServer.Shutdown(ctx)
	}
	return nil
}

// Handlers will return handlers
func (s *Server) Handlers() *nrhttprouter.Router {

	// Create a new router
	s.Router = apirouter.New()

	// Custom logger
	s.Router.Logger = s.Config.Services.Log

	// Turned on all CORs for now
	s.Router.CrossOriginEnabled = true
	s.Router.CrossOriginAllowCredentials = true
	s.Router.CrossOriginAllowOrigin = wildcard
	s.Router.CrossOriginAllowHeaders = wildcard
	s.Router.CrossOriginAllowMethods = strings.Join([]string{
		http.MethodDelete,
		http.MethodGet,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	}, ",")

	// Register all actions (routes / handlers)
	base.RegisterRoutes(s.Router, s.Config, s.P2pServer)

	// Return the router
	return s.Router.HTTPRouter
}
