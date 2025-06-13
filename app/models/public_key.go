package models

import (
	"context"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/bitcoin-sv/alert-system/utils"
	"github.com/mrz1836/go-datastore"
)

// PublicKey is an object representing a public key
type PublicKey struct {
	// Base model
	model.Model `bson:",inline"`

	// Model specific fields
	ID             uint64 `json:"id" toml:"id" yaml:"id" bson:"_id" gorm:"primaryKey;comment:This is a unique identifier"`
	Key            string `json:"key" toml:"key" yaml:"key" bson:"key" gorm:"<-;type:char(66);index;comment:This is the key"`
	LastUpdateHash string `json:"last_update_hash" toml:"last_update_hash" yaml:"last_update_hash" bson:"last_update_hash" gorm:"<-;type:char(64);index;comment:This is the last update hash"`
	Active         bool   `json:"active" toml:"active" yaml:"active" bson:"active" gorm:"<-;type:boolean;index;comment:This is the active flag"`
}

// NewPublicKey creates a new public key
func NewPublicKey(opts ...model.Options) *PublicKey {
	return &PublicKey{
		Model: *model.NewBaseModel(model.NamePublicKey, opts...),
	}
}

// Name will get the name of the model
func (m *PublicKey) Name() string {
	return model.NamePublicKey.String()
}

// GetTableName will get the database table name of the model
func (m *PublicKey) GetTableName() string {
	return model.TablePublicKeys
}

// GetID will get the model ID
func (m *PublicKey) GetID() uint64 {
	return m.ID
}

// Display filter the model for display
func (m *PublicKey) Display() interface{} {
	return m
}

// Migrate will run model-specific migrations on startup
func (m *PublicKey) Migrate(client datastore.ClientInterface) error {
	return client.IndexMetadata(client.GetTableName(model.TablePublicKeys), model.MetadataField)
}

// BeginSaveWithTx will start saving the model into the Datastore with the provided transaction
func (m *PublicKey) BeginSaveWithTx(ctx context.Context, tx *datastore.Transaction) ([]model.BaseInterface, error) {
	return model.BeginSaveWithTx(ctx, tx, m)
}

// Save will save the model into the Datastore
func (m *PublicKey) Save(ctx context.Context) error {
	return model.Save(ctx, m)
}

// GetActivePublicKey will get the active public key
func GetActivePublicKey(ctx context.Context, metadata *model.Metadata, opts ...model.Options) ([]*PublicKey, error) {

	// Set the conditions
	conditions := &map[string]interface{}{
		utils.FieldActive: true, // Active flag is true
		utils.FieldDeletedAt: map[string]interface{}{ // IS NULL
			utils.ExistsCondition: false,
		},
	}

	// Set the query params
	queryParams := &datastore.QueryParams{
		Page:          1,
		PageSize:      10,
		OrderByField:  utils.FieldID,
		SortDirection: utils.SortAscending,
	}

	// Get the records
	modelItems := make([]*PublicKey, 0)
	if err := model.GetModelsByConditions(
		ctx, model.NamePublicKey, &modelItems, metadata, conditions, queryParams, opts...,
	); err != nil {
		return nil, err
	}

	return modelItems, nil
}

// ClearActivePublicKeys will clear the active public keys
// todo this needs to be refactored to use model update/save
func ClearActivePublicKeys(_ context.Context, ds datastore.ClientInterface) error {
	// Execute the query
	tx := ds.Execute("").Exec(
		"UPDATE "+config.DatabasePrefix+"_"+model.TablePublicKeys+" SET "+utils.FieldActive+" = ?",
		false,
	).Begin()

	// Commit the transaction
	tx.Commit()

	// Return the error
	return tx.Error
}
