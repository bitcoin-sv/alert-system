package model

import (
	"github.com/bitcoin-sv/alert-system/app/config"
)

// Options allow functional options to be supplied
// that overwrite default model options
type Options func(m *Model)

// New set this model to a new record
func New() Options {
	return func(m *Model) {
		m.New()
	}
}

// WithLogger will set the current logger on the model
func WithLogger(loggerClient config.LoggerInterface) Options {
	return func(m *Model) {
		if loggerClient != nil {
			m.logger = loggerClient
		}
	}
}

// WithAllDependencies will set all found dependencies on the model
func WithAllDependencies(dependencies *config.Config) Options {
	return func(m *Model) {
		if dependencies != nil {
			m.dependencies = dependencies
		}
	}
}

// WithDebug will set the debugging flag
func WithDebug() Options {
	return func(m *Model) {
		m.debug = true
	}
}

// WithMetadata will add the metadata record to the model
func WithMetadata(key string, value interface{}) Options {
	return func(m *Model) {
		if m.Metadata == nil {
			m.Metadata = make(Metadata)
		}
		m.Metadata.SetKey(key, value)
	}
}

// WithMetadatas will add multiple metadata records to the model
func WithMetadatas(metadata map[string]interface{}) Options {
	return func(m *Model) {
		if len(metadata) > 0 {
			if m.Metadata == nil {
				m.Metadata = make(Metadata)
			}
			for key, value := range metadata {
				m.Metadata.SetKey(key, value)
			}
		}
	}
}
