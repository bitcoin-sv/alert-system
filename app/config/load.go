package config

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mrz1836/go-datastore"
	"github.com/ordishs/gocore"
)

// LoadConfig will load the configuration and services
// models is a list of models to auto-migrate when the datastore is created
func LoadConfig(ctx context.Context, models []interface{}, isTesting bool) (_appConfig *Config, err error) {
	var ok bool

	// Load the database path
	dbPath, _ := gocore.Config().Get("ALERT_SYSTEM_DATABASE_PATH")
	if dbPath == "" {
		dbPath = DatabasePathDefault
	}
	// Sync the configuration struct
	_appConfig = &Config{
		RequestLogging: true,
		Services:       &Services{},
		Datastore: &DatastoreConfig{
			AutoMigrate: true,
			Engine:      datastore.SQLite,
			TablePrefix: DatabasePrefix,
			Debug:       false,
			SQLite: &datastore.SQLiteConfig{
				CommonConfig: datastore.CommonConfig{
					Debug:              false,
					MaxIdleConnections: 1,
					MaxOpenConnections: 1,
					TablePrefix:        DatabasePrefix,
				},
				Shared:       false,
				DatabasePath: dbPath,
			},
		},
		WebServer: &WebServerConfig{
			IdleTimeout:  60 * time.Second, // For idle connections
			Port:         "3000",           // Default port
			ReadTimeout:  15 * time.Second, // For reading the request
			WriteTimeout: 15 * time.Second, // For writing the response
		},
	}

	// Load the logger service (gocore.Logger meets the LoggerInterface)
	_appConfig.Services.Log = &ExtendedLogger{
		Logger: gocore.Log(ApplicationName),
	}

	// Load the RPC user
	if _appConfig.RPCUser, ok = gocore.Config().Get("RPC_USER"); !ok {
		return nil, ErrNoRPCUser
	}

	// Load the RPC password
	if _appConfig.RPCPassword, ok = gocore.Config().Get("RPC_PASSWORD"); !ok {
		return nil, ErrNoRPCPassword
	}

	// Load the RPC host
	if _appConfig.RPCHost, ok = gocore.Config().Get("RPC_HOST"); !ok {
		return nil, ErrNoRPCHost
	}

	// Load the P2P Bootstrap peer
	if _appConfig.P2PBootstrapPeer, ok = gocore.Config().Get("P2P_BOOTSTRAP_PEER"); !ok {
		_appConfig.P2PBootstrapPeer = ""
	}

	// Load the P2P alert system protocol ID
	if _appConfig.P2PAlertSystemProtocolID, ok = gocore.Config().Get("P2P_ALERT_SYSTEM_PROTOCOL_ID"); !ok {
		_appConfig.P2PAlertSystemProtocolID = DefaultAlertSystemProtocolID
	}

	// Load the private key path
	// If not found, create a default one
	if _appConfig.P2PPrivateKeyPath, ok = gocore.Config().Get("P2P_PRIVATE_KEY_PATH"); !ok {
		if err = _appConfig.createPrivateKeyDirectory(); err != nil {
			return nil, err
		}
	}

	// Load the p2p ip
	if _appConfig.P2PIP, ok = gocore.Config().Get("P2P_IP"); !ok {
		return nil, ErrNoP2PIP
	}

	// Load the p2p port
	if _appConfig.P2PPort, ok = gocore.Config().Get("P2P_PORT"); !ok {
		return nil, ErrNoP2PPort
	}

	// Load the webhook URL (if set - this is optional)
	if _appConfig.AlertWebhookURL, ok = gocore.Config().Get("ALERT_WEBHOOK_URL"); !ok {
		_appConfig.Services.Log.Debugf("webhook url is not configured, webhook usage is disabled")
	}

	// Set the node config (either a real node or a mock node)
	if !isTesting {
		_appConfig.Services.Node = NewNodeConfig(_appConfig.RPCUser, _appConfig.RPCPassword, _appConfig.RPCHost)
	} else {
		_appConfig.Services.Node = NewNodeMock(_appConfig.RPCUser, _appConfig.RPCPassword, _appConfig.RPCHost)
	}

	// Use sql in-memory for testing
	// todo this could come from a test struct or test env file
	if isTesting {
		_appConfig.Datastore.AutoMigrate = true
		_appConfig.Datastore.Engine = datastore.SQLite
		_appConfig.Datastore.TablePrefix = DatabasePrefix
		_appConfig.Datastore.SQLite = &datastore.SQLiteConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:              true,
				MaxIdleConnections: 1,
				MaxOpenConnections: 1,
			},
			Shared:       false,
			DatabasePath: "",
		}
	}

	// Load an HTTP client
	_appConfig.Services.HTTPClient = http.DefaultClient

	// Load the datastore service
	if err = _appConfig.loadDatastore(ctx, models); err != nil {
		return nil, err
	}

	return
}

// createPrivateKeyDirectory will create the private key directory
func (c *Config) createPrivateKeyDirectory() error {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to initialize p2p private key file: %w", err)
	}
	if err = os.Mkdir(fmt.Sprintf("%s/%s", dirName, LocalPrivateKeyDirectory), 0750); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("failed to ensure %s dir exists: %w", LocalPrivateKeyDirectory, err)
	}
	c.P2PPrivateKeyPath = fmt.Sprintf("%s/%s/%s", dirName, LocalPrivateKeyDirectory, LocalPrivateKeyDefault)
	return nil
}

// CloseAll will close all connections to all services
func (c *Config) CloseAll(ctx context.Context) {

	// No services to close
	if c.Services == nil {
		return
	}

	// Close the datastore
	if c.Services.Datastore != nil {
		_ = c.Services.Datastore.Close(ctx)
		c.Services.Datastore = nil
	}
}
