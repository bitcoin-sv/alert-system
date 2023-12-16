// Package config provides configuration for the application
package config

import (
	"net/http"
	"time"

	"github.com/mrz1836/go-datastore"
)

// Application configuration constants
var (
	ApplicationName              = "alert_system"                                                                         // Application name used in places where we need an application name space
	DatabasePathDefault          = "alert_system_datastore.db"                                                            // Default database path (Sqlite)
	DatabasePrefix               = "alert_system"                                                                         // Default database prefix
	DefaultServerShutdown        = 5 * time.Second                                                                        // Default server shutdown delay time (to finish any requests or internal processes)
	LocalPrivateKeyDefault       = "alert_system_private_key"                                                             // Default local private key
	LocalPrivateKeyDirectory     = ".bitcoin"                                                                             // Default local private key directory
	SeedIpfsNode                 = "/ip4/68.183.57.231/tcp/9906/p2p/12D3KooWQs6ptKvoKNHurCzqRaVp3uFs9731NQwS3AmVcNc2TGpb" // Default seed IPFS node
	DefaultAlertSystemProtocolID = "/bitcoin/alert-system/1.0.1"                                                          // Default alert system protocol for libp2p syncing
)

// The global configuration settings
type (

	// Config is the global configuration settings
	Config struct {
		AlertWebhookURL          string           `json:"alert_webhook_url"`            // AlertWebhookURL is the URL for the alert webhook
		Datastore                *DatastoreConfig `json:"datastore"`                    // Datastore's configuration
		P2PIP                    string           `json:"p2p_ip"`                       // P2PIP is the IP address for the P2P server
		P2PPort                  string           `json:"p2p_port"`                     // P2PPort is the port for the P2P server
		P2PPrivateKeyPath        string           `json:"p2p_private_key_path"`         // P2PPrivateKeyPath is the path to the private key
		P2PBootstrapPeer         string           `json:"p2p_bootstrap_peer"`           // P2PBootstrapPeer is the bootstrap peer for the libp2p network
		P2PAlertSystemProtocolID string           `json:"p2p_alert_system_protocol_id"` // P2PAlertSystemProtocolID is the protocol ID to use on the libp2p network for alert system communication
		RPCHost                  string           `json:"rpc_host"`                     // RPCHost is the RPC host
		RPCPassword              string           `json:"rpc_password"`                 // RPCPassword is the RPC password
		RPCUser                  string           `json:"rpc_user"`                     // RPCUser is the RPC username
		RequestLogging           bool             `json:"request_logging"`              // Toggle for verbose request logging (API requests)
		Services                 *Services        `json:"-"`                            // Services is the global services
		WebServer                *WebServerConfig `json:"web_server"`                   // WebServer is the configuration for the web HTTP Server
	}

	// DatastoreConfig is the configuration for the datastore
	DatastoreConfig struct {
		AutoMigrate bool                    `json:"auto_migrate"` // Loads a blank database
		Debug       bool                    `json:"debug"`        // True for sql statements
		Engine      datastore.Engine        `json:"engine"`       // MySQL, Postgres, SQLite
		Password    string                  `json:"password"`     // Used for MySQL or Postgresql
		SQLite      *datastore.SQLiteConfig `json:"sqlite"`       // Configuration for SQLite
		SQLRead     *datastore.SQLConfig    `json:"sql_read"`     // Configuration for MySQL or Postgres
		SQLWrite    *datastore.SQLConfig    `json:"sql_write"`    // Configuration for MySQL or Postgres
		TablePrefix string                  `json:"table_prefix"` // pre_table_name (pre)
	}

	// HTTPInterface is used for the HTTP client
	HTTPInterface interface {
		Do(req *http.Request) (*http.Response, error)
	}

	// Node is the configuration and functions for interacting with a node
	Node struct {
		RPCHost     string `json:"rpc_host"`     // RPCHost is the RPC host
		RPCPassword string `json:"rpc_password"` // RPCPassword is the RPC password
		RPCUser     string `json:"rpc_user"`     // RPCUser is the RPC username
	}

	// WebServerConfig is a configuration for the web HTTP Server
	WebServerConfig struct {
		IdleTimeout  time.Duration `json:"idle_timeout"`  // 60s
		Port         string        `json:"port"`          // 3000
		ReadTimeout  time.Duration `json:"read_timeout"`  // 15s
		WriteTimeout time.Duration `json:"write_timeout"` // 15s
	}

	// Services is the global services
	Services struct {
		Datastore  datastore.ClientInterface // Datastore interface
		Log        LoggerInterface           // Logger interface
		Node       NodeInterface             // Node interface
		HTTPClient HTTPInterface             // HTTP client interface
	}
)
