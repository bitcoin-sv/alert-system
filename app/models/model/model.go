// Package model provides the base model for all models
package model

import (
	"context"
	"time"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/mrz1836/go-datastore"
	customTypes "github.com/mrz1836/go-datastore/custom_types"
	"github.com/ordishs/gocore"
)

// Model is the generic model field(s) and interface(s)
//
// gorm: https://gorm.io/docs/models.html
type Model struct {
	// ID string  `json:"id" toml:"id" yaml:"id" gorm:"primaryKey"`  (@mrz: custom per table)

	CreatedAt time.Time            `json:"created_at" toml:"created_at" yaml:"created_at" bson:"created_at" gorm:"comment:The time that the record was originally created"`
	DeletedAt customTypes.NullTime `json:"deleted_at" toml:"deleted_at" yaml:"deleted_at" bson:"deleted_at,omitempty" gorm:"index;comment:The time the record was marked as deleted"`
	Metadata  Metadata             `toml:"metadata,omitempty" yaml:"metadata,omitempty" bson:"metadata,omitempty" json:"metadata,omitempty" gorm:"type:json;comment:The JSON metadata for the record"`
	UpdatedAt time.Time            `json:"updated_at" toml:"updated_at" yaml:"updated_at" bson:"updated_at,omitempty" gorm:"comment:The time that the record was last updated"`

	// Private fields
	debug        bool                   // Set from the parent config if debugging is turned on/off
	dependencies *config.Config         // Application dependencies (app config, services, datastore, cachestore, etc)
	logger       config.LoggerInterface // Internal logging
	name         Name                   // Name of model (table name)
	newRecord    bool                   // Determine if the record is new (create vs update)
}

// BaseInterface is the interface that all models share
type BaseInterface interface {
	AfterCreated(ctx context.Context) error
	AfterDeleted(ctx context.Context) error
	AfterUpdated(ctx context.Context) error
	BeforeCreating(ctx context.Context) error
	BeforeUpdating(ctx context.Context) error
	BeginSaveWithTx(ctx context.Context, tx *datastore.Transaction) ([]BaseInterface, error)
	ChildModels() []BaseInterface
	Config() *config.Config
	Datastore() datastore.ClientInterface
	Debug(enable bool)
	DebugLog(ctx context.Context, text string)
	Display() interface{}
	GetID() uint64
	GetOptions(isNewRecord bool) []Options
	GetTableName() string
	IsNew() bool
	Logger() config.LoggerInterface
	Migrate(client datastore.ClientInterface) error
	Name() string
	New()
	NotNew()
	Save(ctx context.Context) error
	SetMetaData(key string, value interface{})
	SetOptions(opts ...Options)
	SetRecordTime(isNew bool)
	UpdateMetadata(metadata Metadata)
}

// Name is the model name type
type Name string

// NewBaseModel create an empty base model
func NewBaseModel(name Name, opts ...Options) (m *Model) {
	m = &Model{name: name}
	m.SetOptions(opts...)

	// Set default logger IF NOT SET via options
	if m.logger == nil {
		m.logger = &config.ExtendedLogger{
			Logger: gocore.Log(config.ApplicationName),
		}
	}

	return
}

/*
// DisplayModels process the (slice) of model(s) for display
func DisplayModels(models interface{}) interface{} {
	if models == nil {
		return nil
	}

	s := reflect.ValueOf(models)
	if s.IsNil() {
		return nil
	}
	if s.Kind() == reflect.Slice {
		for i := 0; i < s.Len(); i++ {
			s.Index(i).MethodByName("Display").Call([]reflect.Value{})
		}
	} else {
		s.MethodByName("Display").Call([]reflect.Value{})
	}

	return models
}
*/

// String is the string version of the name
func (n Name) String() string {
	return string(n)
}

// IsEmpty tests if the model name is empty
func (n Name) IsEmpty() bool {
	return n == NameEmpty || n == ""
}
