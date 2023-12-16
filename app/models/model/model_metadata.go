package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/mrz1836/go-datastore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	// MetadataField is the field name used for metadata (params)
	MetadataField = "metadata"
)

// Metadata is an object representing the metadata about the related record (standard across all tables)
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type Metadata map[string]interface{}

// GormDataType type in gorm
func (m Metadata) GormDataType() string {
	return gormTypeText
}

// SetKey will set a key/value pair in the metadata
func (m *Metadata) SetKey(key string, value interface{}) {
	if *m == nil { // Extra overhead here for this check, but it now works if Metadata is nil
		*m = make(Metadata)
	}
	(*m)[key] = value
}

// GetKey will get a key from the metadata
func (m *Metadata) GetKey(key string) interface{} {
	return (*m)[key]
}

// Scan scan value into Json, implements sql.Scanner interface
func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	xType := fmt.Sprintf("%T", value)
	var byteValue []byte
	if xType == "string" {
		byteValue = []byte(value.(string))
	} else {
		byteValue = value.([]byte)
	}
	if bytes.Equal(byteValue, []byte("")) || bytes.Equal(byteValue, []byte("\"\"")) {
		return nil
	}

	return json.Unmarshal(byteValue, &m)
}

// Value return json value, implement driver.Valuer interface
func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return string(marshal), nil
}

// GormDBDataType the gorm data type for metadata
func (Metadata) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	if db.Dialector.Name() == datastore.Postgres {
		return datastore.JSONB
	}
	return datastore.JSON
}

// MarshalBSONValue method is called by bson.Marshal in Mongo for type = Metadata
func (m *Metadata) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if m == nil || len(*m) == 0 {
		return bson.TypeNull, nil, nil
	}

	metadata := make([]map[string]interface{}, 0)
	for key, value := range *m {
		metadata = append(metadata, map[string]interface{}{
			"k": key,
			"v": value,
		})
	}

	return bson.MarshalValue(metadata)
}

// UnmarshalBSONValue method is called by bson.Unmarshal in Mongo for type = Metadata
func (m *Metadata) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	raw := bson.RawValue{Type: t, Value: data}

	if raw.Value == nil {
		return nil
	}

	var uMap []map[string]interface{}
	if err := raw.Unmarshal(&uMap); err != nil {
		return err
	}

	*m = make(Metadata)
	for _, meta := range uMap {
		key := meta["k"].(string)
		(*m)[key] = meta["v"]
	}

	return nil
}

// MarshalMetadata will marshal the custom type
func MarshalMetadata(m Metadata) graphql.Marshaler {
	if m == nil {
		return graphql.Null
	}
	return graphql.MarshalMap(m)
}

// UnmarshalMetadata will unmarshal the custom type
func UnmarshalMetadata(v interface{}) (Metadata, error) {
	if v == nil {
		return nil, nil
	}
	return graphql.UnmarshalMap(v)
}
