// Package main is the entry point for the alert-system
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/bitcoin-sv/alert-system/app/p2p"
	"github.com/bitcoin-sv/alert-system/app/webserver"
)

// main is the entry point for the alert-system
func main() {

	// Load the configuration and services
	_appConfig, err := config.LoadDependencies(context.Background(), models.BaseModels, false)
	if err != nil {
		log.Fatalf("error loading configuration: %s", err.Error())
	}
	defer func() {
		_appConfig.CloseAll(context.Background())
	}()

	// Ensure we have the genesis alert in the database
	if err = models.CreateGenesisAlert(
		context.Background(), model.WithAllDependencies(_appConfig),
	); err != nil {
		_appConfig.Services.Log.Fatalf("error creating genesis alert: %s", err.Error())
	}

	// Ensure that RPC connection is valid
	if !_appConfig.DisableRPCVerification {
		if _, err = _appConfig.Services.Node.BestBlockHash(context.Background()); err != nil {
			_appConfig.Services.Log.Errorf("error talking to Bitcoin node with supplied RPC credentials: %s", err.Error())
			return
		}
	}

	// Create the p2p server
	var p2pServer *p2p.Server
	if p2pServer, err = p2p.NewServer(p2p.ServerOptions{
		TopicNames: []string{_appConfig.P2P.TopicName},
		Config:     _appConfig,
	}); err != nil {
		_appConfig.Services.Log.Fatalf("error creating p2p server: %s", err.Error())
	}

	// Create a new (web) server
	webServer := webserver.NewServer(_appConfig, p2pServer)

	ctx, cancelFunc := context.WithCancel(context.Background())
	// Start the p2p server
	if err = p2pServer.Start(ctx); err != nil {
		_appConfig.Services.Log.Fatalf("error starting p2p server: %s", err.Error())
	}
	// Sync a channel to listen for interrupts
	idleConnectionsClosed := make(chan struct{})
	go func(appConfig *config.Config) {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)

		// Log when a signal is received
		appConfig.Services.Log.Debugf("waiting for interrupt signal")
		<-sigint

		// Log that we are starting the shutdown process
		appConfig.Services.Log.Infof("interrupt signal received, starting shutdown process")

		// We received an interrupt signal, shut down the server
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultServerShutdown)
		defer cancel()
		if err = webServer.Shutdown(ctxTimeout); err != nil {
			appConfig.Services.Log.Infof("error shutting down webserver: %s", err.Error())
		}

		// Shutdown the p2p server
		if err = p2pServer.Stop(ctxTimeout); err != nil {
			appConfig.Services.Log.Infof("error shutting down p2p server: %s", err.Error())
		}
		cancelFunc()
		appConfig.Services.Log.Infof("successfully shut down server")
		close(idleConnectionsClosed)
		if err = appConfig.Services.Log.CloseWriter(); err != nil {
			log.Printf("error closing logger: %s", err)
		}
	}(_appConfig)

	// Serve the web server and then wait endlessly
	webServer.Serve()

	// Wait for the idle connection to close
	<-idleConnectionsClosed
}
