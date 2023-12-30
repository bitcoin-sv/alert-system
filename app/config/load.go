package config

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/mrz1836/go-datastore"
	"github.com/ordishs/gocore"
	"github.com/spf13/viper"
)

// Added a mutex lock for a race-condition
var viperLock sync.Mutex

// isValidEnvironment will return true if the testEnv is a known valid environment
func isValidEnvironment(testEnv string) bool {
	testEnv = strings.ToLower(testEnv)
	for _, env := range environments {
		if env == testEnv {
			return true
		}
	}
	return false
}

// LoadDependencies will load the configuration and services
// models is a list of models to auto-migrate when the datastore is created
// if testing is true, the node will be mocked
func LoadDependencies(ctx context.Context, models []interface{}, isTesting bool) (_appConfig *Config, err error) {

	// Load the config file
	_appConfig, err = LoadConfigFile()
	if err != nil {
		return nil, err
	}

	// Require at least one RPC connection
	if len(_appConfig.RPCConnections) == 0 {
		return nil, ErrNoRPCConnections
	}

	// Ensure the P2P configuration is valid
	if err = requireP2P(_appConfig); err != nil {
		return nil, err
	}

	// Set the node config (either a real node or a mock node)
	if !isTesting {
		// todo support multiple nodes (this is an example)
		for i := range _appConfig.RPCConnections {
			_appConfig.Services.Node = NewNodeConfig(
				_appConfig.RPCConnections[i].User,
				_appConfig.RPCConnections[i].Password,
				_appConfig.RPCConnections[i].Host,
			)
		}
	} else {
		for i := range _appConfig.RPCConnections {
			_appConfig.Services.Node = NewNodeMock(
				_appConfig.RPCConnections[i].User,
				_appConfig.RPCConnections[i].Password,
				_appConfig.RPCConnections[i].Host,
			)
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

// requireP2P will ensure the P2P configuration is valid
func requireP2P(_appConfig *Config) error {

	// Set the P2P alert system protocol ID if it's missing
	if len(_appConfig.P2P.AlertSystemProtocolID) == 0 {
		_appConfig.P2P.AlertSystemProtocolID = DefaultAlertSystemProtocolID
	}

	// Load the private key path
	// If not found, create a default one
	if len(_appConfig.P2P.PrivateKeyPath) == 0 {
		if err := _appConfig.createPrivateKeyDirectory(); err != nil {
			return err
		}
	}

	// Load the p2p ip (local, ip address or domain name)
	// todo better validation of what is a valid IP, domain name or local address
	if len(_appConfig.P2P.IP) < 5 {
		return ErrNoP2PIP
	}

	// Load the p2p port ( >= XX)
	if len(_appConfig.P2P.Port) < 2 {
		return ErrNoP2PPort
	}

	return nil
}

// LoadConfigFile will load the config file and environment variables
func LoadConfigFile() (_appConfig *Config, err error) {

	// Start the configuration struct
	_appConfig = &Config{
		Datastore: DatastoreConfig{
			SQLite:   &datastore.SQLiteConfig{},
			SQLRead:  &datastore.SQLConfig{},
			SQLWrite: &datastore.SQLConfig{},
		},
		P2P:            P2PConfig{},
		Services:       Services{},
		WebServer:      WebServerConfig{},
		RPCConnections: make([]RPCConfig, 0),
	}

	// Check the environment we are running
	environment := os.Getenv(EnvironmentKey)
	if !isValidEnvironment(environment) {
		err = ErrInvalidEnvironment
		return nil, err
	}

	// Lock viper
	viperLock.Lock()

	// Unlock the viper mutex
	defer viperLock.Unlock()

	// Set a replacer for replacing double underscore with nested period
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)

	// Set the prefix
	viper.SetEnvPrefix(EnvironmentPrefix)

	// Use env vars
	viper.AutomaticEnv()

	// Get the embedded envs directory
	var files []fs.DirEntry
	if files, err = envDir.ReadDir("envs"); err != nil {
		return nil, err
	}

	// Set the configuration type
	viper.SetConfigType("json")

	// Do we have a custom config file? (use this instead of the environment file)
	customConfigFileWithPath := os.Getenv(EnvironmentCustomFilePath)
	if len(customConfigFileWithPath) > 0 {
		var b []byte

		// Read the file
		if b, err = os.ReadFile(customConfigFileWithPath); err != nil { //nolint:gosec // This is a custom file path
			return nil, err
		}

		// Read the config
		if err = viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
			return nil, err
		}
	} else {
		// Loop through the various environment files
		for _, file := range files {
			if file.Name() == environment+".json" {
				var f fs.File
				if f, err = envDir.Open("envs/" + file.Name()); err != nil {
					return nil, err
				}
				if err = viper.ReadConfig(f); err != nil {
					return nil, err
				}
			}
		}
	}

	// Unmarshal into values struct
	if err = viper.Unmarshal(&_appConfig); err != nil {
		err = fmt.Errorf("error loading viper values: %w", err)
		return nil, err
	}

	// Load the logger service (gocore.Logger meets the LoggerInterface)
	_appConfig.Services.Log = &ExtendedLogger{
		Logger: gocore.Log(ApplicationName),
	}

	// Log the configuration that was detected and where it was loaded from
	_appConfig.Services.Log.Debug("loaded configuration from: " + viper.ConfigFileUsed())

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
	c.P2P.PrivateKeyPath = fmt.Sprintf("%s/%s/%s", dirName, LocalPrivateKeyDirectory, LocalPrivateKeyDefault)
	return nil
}

// CloseAll will close all connections to all services
func (c *Config) CloseAll(ctx context.Context) {

	// Close the datastore
	if c.Services.Datastore != nil {
		_ = c.Services.Datastore.Close(ctx)
		c.Services.Datastore = nil
	}
}
