package config

import (
	"context"

	"github.com/mrz1836/go-logger"

	"github.com/mrz1836/go-datastore"
)

// loadDatastore will load an instance of Datastore into the dependencies
func (c *Config) loadDatastore(ctx context.Context, models []interface{}) error {

	// Sync collecting the options
	var options []datastore.ClientOps
	//TODO: pass in our own logger, but for now this doesn't work so i'm just going to silently log for now
	options = append(options, datastore.WithLogger(logger.NewGormLogger(false, 0)))
	// Select the datastore
	switch c.Datastore.Engine {
	case datastore.SQLite:
		options = append(options, datastore.WithSQLite(&datastore.SQLiteConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:              c.Datastore.Debug,
				MaxIdleConnections: c.Datastore.SQLite.MaxIdleConnections,
				MaxOpenConnections: c.Datastore.SQLite.MaxOpenConnections,
				TablePrefix:        c.Datastore.TablePrefix,
			},
			DatabasePath: c.Datastore.SQLite.DatabasePath, // "" for in memory
			Shared:       c.Datastore.SQLite.Shared,
		}))
	case datastore.MySQL, datastore.PostgreSQL:

		// Set the pw if not set
		if len(c.Datastore.Password) > 0 && len(c.Datastore.SQLRead.Password) == 0 {
			c.Datastore.SQLRead.Password = c.Datastore.Password
		}
		if len(c.Datastore.Password) > 0 && len(c.Datastore.SQLWrite.Password) == 0 {
			c.Datastore.SQLWrite.Password = c.Datastore.Password
		}

		// Create the read/write options
		options = append(options, datastore.WithSQL(c.Datastore.Engine, []*datastore.SQLConfig{
			{ // MASTER - WRITE
				CommonConfig: datastore.CommonConfig{
					Debug:                 c.Datastore.Debug,
					MaxConnectionIdleTime: c.Datastore.SQLWrite.MaxConnectionTime,
					MaxConnectionTime:     c.Datastore.SQLWrite.MaxConnectionTime,
					MaxIdleConnections:    c.Datastore.SQLWrite.MaxIdleConnections,
					MaxOpenConnections:    c.Datastore.SQLWrite.MaxOpenConnections,
					TablePrefix:           c.Datastore.TablePrefix,
				},
				Driver:    c.Datastore.Engine.String(),
				Host:      c.Datastore.SQLWrite.Host,
				Name:      c.Datastore.SQLWrite.Name,
				Password:  c.Datastore.SQLWrite.Password,
				Port:      c.Datastore.SQLWrite.Port,
				Replica:   false,
				TimeZone:  c.Datastore.SQLWrite.TimeZone,
				TxTimeout: c.Datastore.SQLWrite.TxTimeout,
				User:      c.Datastore.SQLWrite.User,
				SslMode:   c.Datastore.SQLWrite.SslMode,
			},
			{ // READ REPLICA
				CommonConfig: datastore.CommonConfig{
					Debug:                 c.Datastore.Debug,
					MaxConnectionIdleTime: c.Datastore.SQLRead.MaxConnectionTime,
					MaxConnectionTime:     c.Datastore.SQLRead.MaxConnectionTime,
					MaxIdleConnections:    c.Datastore.SQLRead.MaxIdleConnections,
					MaxOpenConnections:    c.Datastore.SQLRead.MaxOpenConnections,
					TablePrefix:           c.Datastore.TablePrefix,
				},
				Driver:    c.Datastore.Engine.String(),
				Host:      c.Datastore.SQLRead.Host,
				Name:      c.Datastore.SQLRead.Name,
				Password:  c.Datastore.SQLRead.Password,
				Port:      c.Datastore.SQLRead.Port,
				Replica:   true,
				TimeZone:  c.Datastore.SQLRead.TimeZone,
				TxTimeout: c.Datastore.SQLRead.TxTimeout,
				User:      c.Datastore.SQLRead.User,
				SslMode:   c.Datastore.SQLRead.SslMode,
			},
		}))
	case datastore.Empty, datastore.MongoDB:
		return ErrDatastoreUnsupported
	default:
		return ErrDatastoreUnsupported
	}

	// Add the auto migration option if enabled
	if c.Datastore.AutoMigrate && models != nil {
		options = append(options, datastore.WithAutoMigrate(models...))
	}

	// Load datastore or return an error
	var err error
	c.Services.Datastore, err = datastore.NewClient(ctx, options...)
	return err
}
