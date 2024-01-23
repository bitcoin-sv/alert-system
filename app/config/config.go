// Package config provides configuration for the application
package config

import (
	"embed"
	"net/http"
	"time"

	"github.com/mrz1836/go-datastore"
)

//go:embed envs
var envDir embed.FS // This is used for the config files

// Constants for the environment
const (
	EnvironmentCustomFilePath = "ALERT_SYSTEM_CONFIG_FILEPATH" // Environment variable key for custom config file path
	EnvironmentKey            = "ALERT_SYSTEM_ENVIRONMENT"     // Environment variable key
	EnvironmentLocal          = "local"                        // Environment for local development
	EnvironmentPrefix         = "alert_system"                 // Prefix for all environment variables
	EnvironmentProduction     = "production"                   // Environment for production
	EnvironmentTest           = "test"                         // Environment for testing
	EnvironmentStn            = "stn"                          // Environment for STN testing
)

// Local variables for configuration
var (
	environments = []interface{}{
		EnvironmentLocal,
		EnvironmentProduction,
		EnvironmentTest,
		EnvironmentStn,
	}
)

// Application configuration constants
var (
	ApplicationName              = "alert_system"                // Application name used in places where we need an application name space
	DatabasePrefix               = "alert_system"                // Default database prefix
	DefaultAlertSystemProtocolID = "/bitcoin/alert-system/0.0.1" // Default alert system protocol for libp2p syncing
	DefaultTopicName             = "alert_system"                // Default alert system topic name for libp2p subscription
	DefaultServerShutdown        = 5 * time.Second               // Default server shutdown delay time (to finish any requests or internal processes)
	LocalPrivateKeyDefault       = "alert_system_private_key"    // Default local private key
	LocalPrivateKeyDirectory     = ".bitcoin"                    // Default local private key directory
)

// The global configuration settings
type (

	// Config is the global configuration settings
	Config struct {
		AlertWebhookURL string          `json:"alert_webhook_url" mapstructure:"alert_webhook_url"` // AlertWebhookURL is the URL for the alert webhook
		Datastore       DatastoreConfig `json:"datastore" mapstructure:"datastore"`                 // Datastore's configuration
		P2P             P2PConfig       `json:"p2p" mapstructure:"p2p"`                             // P2P is the configuration for the P2P server
		RPCConnections  []RPCConfig     `json:"rpc_connections" mapstructure:"rpc_connections"`     // RPCConnections is a list of RPC connections
		RequestLogging  bool            `json:"request_logging" mapstructure:"request_logging"`     // Toggle for verbose request logging (API requests)
		Services        Services        `json:"-" mapstructure:"services"`                          // Services is the global services
		WebServer       WebServerConfig `json:"web_server" mapstructure:"web_server"`               // WebServer is the configuration for the web HTTP Server
	}

	// DatastoreConfig is the configuration for the datastore
	DatastoreConfig struct {
		AutoMigrate bool                    `json:"auto_migrate" mapstructure:"auto_migrate"` // Loads a blank database
		Debug       bool                    `json:"debug" mapstructure:"debug"`               // True for sql statements
		Engine      datastore.Engine        `json:"engine" mapstructure:"engine"`             // MySQL, Postgres, SQLite
		Password    string                  `json:"password" mapstructure:"password"`         // Used for MySQL or Postgresql
		SQLite      *datastore.SQLiteConfig `json:"sqlite" mapstructure:"sqlite"`             // Configuration for SQLite
		SQLRead     *datastore.SQLConfig    `json:"sql_read" mapstructure:"sql_read"`         // Configuration for MySQL or Postgres
		SQLWrite    *datastore.SQLConfig    `json:"sql_write" mapstructure:"sql_write"`       // Configuration for MySQL or Postgres
		TablePrefix string                  `json:"table_prefix" mapstructure:"table_prefix"` // pre_table_name (pre)
	}

	// HTTPInterface is used for the HTTP client
	HTTPInterface interface {
		Do(req *http.Request) (*http.Response, error)
	}

	// Node is the configuration and functions for interacting with a node
	Node struct {
		RPCHost     string `json:"rpc_host" mapstructure:"rpc_host"`         // RPCHost is the RPC host
		RPCPassword string `json:"rpc_password" mapstructure:"rpc_password"` // RPCPassword is the RPC password
		RPCUser     string `json:"rpc_user" mapstructure:"rpc_user"`         // RPCUser is the RPC username
	}

	// P2PConfig is the configuration for the P2P server and connection
	P2PConfig struct {
		AlertSystemProtocolID string `json:"alert_system_protocol_id" mapstructure:"alert_system_protocol_id"` // AlertSystemProtocolID is the protocol ID to use on the libp2p network for alert system communication
		BootstrapPeer         string `json:"bootstrap_peer" mapstructure:"bootstrap_peer"`                     // BootstrapPeer is the bootstrap peer for the libp2p network
		IP                    string `json:"ip" mapstructure:"ip"`                                             // IP is the IP address for the P2P server
		Port                  string `json:"port" mapstructure:"port"`                                         // Port is the port for the P2P server
		PrivateKeyPath        string `json:"private_key_path" mapstructure:"private_key_path"`                 // PrivateKeyPath is the path to the private key
		TopicName             string `json:"topic_name" mapstructure:"topic_name"`                             // TopicName is the name of the topic to subscribe to
	}

	// RPCConfig is the configuration for the RPC client
	RPCConfig struct {
		Host     string `json:"host" mapstructure:"host"`         // Host is the RPC host
		Password string `json:"password" mapstructure:"password"` // Password is the RPC password
		User     string `json:"user" mapstructure:"user"`         // User is the RPC username
	}

	// Services is the global services
	Services struct {
		Datastore  datastore.ClientInterface // Datastore interface
		Log        LoggerInterface           // Logger interface
		Node       NodeInterface             // Node interface
		HTTPClient HTTPInterface             // HTTP client interface
	}

	// WebServerConfig is a configuration for the web HTTP Server
	WebServerConfig struct {
		IdleTimeout  time.Duration `json:"idle_timeout" mapstructure:"idle_timeout"`   // 60s
		Port         string        `json:"port" mapstructure:"port"`                   // 3000
		ReadTimeout  time.Duration `json:"read_timeout" mapstructure:"read_timeout"`   // 15s
		WriteTimeout time.Duration `json:"write_timeout" mapstructure:"write_timeout"` // 15s
	}
)
