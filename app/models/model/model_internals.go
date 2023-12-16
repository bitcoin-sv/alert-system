package model

import (
	"context"
	"time"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/mrz1836/go-datastore"
)

// AfterCreated will fire after the model is created in the Datastore
func (m *Model) AfterCreated(_ context.Context) error {
	// m.DebugLog(ctx, "starting: "+m.Name()+" AfterCreated hook...")
	// m.DebugLog(ctx, "end: "+m.Name()+" AfterCreated hook")
	return nil
}

// AfterDeleted will fire after a successful delete in the Datastore
func (m *Model) AfterDeleted(_ context.Context) error {
	// m.DebugLog(ctx, "starting: "+m.Name()+" AfterDelete hook...")
	// m.DebugLog(ctx, "end: "+m.Name()+" AfterDelete hook")
	return nil
}

// AfterUpdated will fire after a successful update into the Datastore
func (m *Model) AfterUpdated(_ context.Context) error {
	// m.DebugLog(ctx, "starting: "+m.Name()+" AfterUpdated hook...")
	// m.DebugLog(ctx, "end: "+m.Name()+" AfterUpdated hook")
	return nil
}

// BeforeCreating will fire before creating a model in the Datastore
func (m *Model) BeforeCreating(_ context.Context) error {
	// m.DebugLog(ctx, "starting: "+m.Name()+" BeforeCreate hook...")
	// m.DebugLog(ctx, "end: "+m.Name()+" BeforeCreate hook")
	return nil
}

// BeforeUpdating will fire before updating a model in the Datastore
func (m *Model) BeforeUpdating(_ context.Context) error {
	// m.DebugLog(ctx, "starting: "+m.Name()+" BeforeUpdate hook...")
	// m.DebugLog(ctx, "end: "+m.Name()+" BeforeUpdate hook")
	return nil
}

// ChildModels will return any child models
func (m *Model) ChildModels() []BaseInterface {
	return nil
}

// Config will return the current configuration
func (m *Model) Config() *config.Config {
	return m.dependencies
}

// Debug will set the debug flag
func (m *Model) Debug(enabled bool) {
	m.debug = enabled
}

// DebugLog will display verbose logs
func (m *Model) DebugLog(ctx context.Context, text string) {
	if m.debug {
		m.Logger().Info(ctx, text)
	}
}

// Datastore will return the current datastore
func (m *Model) Datastore() datastore.ClientInterface {
	return m.dependencies.Services.Datastore
}

// Display filter the model for display
func (m *Model) Display() interface{} {
	return m
}

// enrich is run after getting a record from the database
func (m *Model) enrich(name Name, opts ...Options) {
	// Overwrite defaults
	m.name = name
	m.SetOptions(opts...)
}

// GetID will get the model id, if overwritten in the actual model
func (m *Model) GetID() string {
	return ""
}

// GetTableName will get the table name
func (m *Model) GetTableName() string {
	return TableEmpty
}

// GetOptions will get the options that are set on that model
func (m *Model) GetOptions(isNewRecord bool) (opts []Options) {

	// All dependencies
	if m.dependencies != nil {
		opts = append(opts, WithAllDependencies(m.dependencies))
	}

	// Logger client from the model
	if m.logger != nil {
		opts = append(opts, WithLogger(m.logger))
	}

	// Debugging
	if m.debug {
		opts = append(opts, WithDebug())
	}

	// New record flag
	if isNewRecord {
		opts = append(opts, New())
	}

	return
}

// IsNew returns true if the model is (or was) a new record
func (m *Model) IsNew() bool {
	return m.newRecord
}

// Logger will return the Logger if it exists
func (m *Model) Logger() config.LoggerInterface {
	return m.logger
}

// Migrate will run custom migrations for the model
func (m *Model) Migrate(_ datastore.ClientInterface) error {
	return nil
}

// Name will get the collection name (model)
func (m *Model) Name() string {
	return m.name.String()
}

// New will set the record to new
func (m *Model) New() {
	m.newRecord = true
}

// NotNew sets newRecord to false
func (m *Model) NotNew() {
	m.newRecord = false
}

// BeginSaveWithTx will start saving the model into the Datastore with the provided transaction
func (m *Model) BeginSaveWithTx(_ context.Context, _ *datastore.Transaction) ([]BaseInterface, error) {
	return nil, nil
}

// Save will save the model and child models
func (m *Model) Save(_ context.Context) error {
	return nil
}

// SetOptions will set the options on the model
func (m *Model) SetOptions(opts ...Options) {
	for _, opt := range opts {
		opt(m)
	}
}

// SetRecordTime will set the record timestamps (created is true for a new record)
func (m *Model) SetRecordTime(created bool) {
	if created {
		m.CreatedAt = time.Now().UTC()
		m.UpdatedAt = time.Now().UTC() // Override the default so it's UTC
	} else {
		m.UpdatedAt = time.Now().UTC()
	}
}

// UpdateMetadata will update the metadata on the model
// any key set to nil will be removed, other keys updated or added
func (m *Model) UpdateMetadata(metadata Metadata) {
	if m.Metadata == nil {
		m.Metadata = make(Metadata)
	}

	for key, value := range metadata {
		if value == nil {
			delete(m.Metadata, key)
		} else {
			m.Metadata.SetKey(key, value)
		}
	}
}

// SetMetaData will set the metadata value
func (m *Model) SetMetaData(key string, value interface{}) {
	if m.Metadata == nil {
		m.Metadata = Metadata{}
	}
	m.Metadata.SetKey(key, value)
}
