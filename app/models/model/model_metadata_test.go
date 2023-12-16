package model

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/mrz1836/go-datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestMetaDataScan will test the db Scanner of the Metadata model
func TestMetadata_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		m := Metadata{}
		err := m.Scan(nil)
		require.NoError(t, err)
		assert.Empty(t, len(m))
	})

	t.Run("empty string", func(t *testing.T) {
		m := Metadata{}
		err := m.Scan([]byte("\"\""))
		require.NoError(t, err)
		assert.Empty(t, len(m))
	})

	t.Run("empty string - incorrectly coded", func(t *testing.T) {
		m := Metadata{}
		err := m.Scan([]byte(""))
		require.NoError(t, err)
		assert.Empty(t, len(m))
	})

	t.Run("object", func(t *testing.T) {
		m := Metadata{}
		err := m.Scan([]byte("{\"test\":\"test2\"}"))
		require.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, "test2", m["test"])
	})
}

// TestMetadata_Set will test the db Valuer of the Metadata model
func TestMetadata_Set(t *testing.T) {
	t.Parallel()

	t.Run("set key", func(t *testing.T) {
		m := Metadata{}
		m.SetKey("test", "test2")
		assert.Equal(t, "test2", m.GetKey("test"))
	})
}

// TestMetadata_Get will test the db Valuer of the Metadata model
func TestMetadata_Get(t *testing.T) {
	t.Parallel()

	t.Run("get key", func(t *testing.T) {
		m := Metadata{}
		m.SetKey("test", "test2")
		assert.Equal(t, "test2", m.GetKey("test"))
	})
}

// TestMetadata_Value will test the db Valuer of the Metadata model
func TestMetadata_Value(t *testing.T) {
	t.Parallel()

	t.Run("empty object", func(t *testing.T) {
		m := Metadata{}
		value, err := m.Value()
		require.NoError(t, err)
		assert.Equal(t, "{}", value)
	})

	t.Run("map present", func(t *testing.T) {
		m := Metadata{}
		m["test"] = "test2"
		value, err := m.Value()
		require.NoError(t, err)
		assert.Equal(t, "{\"test\":\"test2\"}", value)
	})
}

// TestMetadata_UnmarshalMetadata will test unmarshalling the metadata
func TestMetadata_UnmarshalMetadata(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		m, err := UnmarshalMetadata(nil)
		require.NoError(t, err)
		assert.Empty(t, len(m))
		assert.IsType(t, Metadata{}, m)
	})

	t.Run("empty string", func(t *testing.T) {
		m, err := UnmarshalMetadata("\"\"")
		require.Error(t, err)
		assert.Empty(t, len(m))
		assert.IsType(t, Metadata{}, m)
	})

	t.Run("empty string - incorrectly coded", func(t *testing.T) {
		m, err := UnmarshalMetadata("")
		require.Error(t, err)
		assert.Empty(t, len(m))
		assert.IsType(t, Metadata{}, m)
	})

	t.Run("object", func(t *testing.T) {
		m, err := UnmarshalMetadata(map[string]interface{}{"test": "test2"})
		require.NoError(t, err)
		assert.Len(t, m, 1)
		assert.IsType(t, Metadata{}, m)
		assert.Equal(t, "test2", m["test"])
	})
}

// TestMetadata_MarshalMetadata will test marshaling the metadata
func TestMetadata_MarshalMetadata(t *testing.T) {
	t.Parallel()

	t.Run("empty object", func(t *testing.T) {
		m := Metadata{}
		writer := MarshalMetadata(m)
		require.NotNil(t, writer)
		b := bytes.NewBufferString("")
		writer.MarshalGQL(b)
		assert.Equal(t, "{}\n", b.String())
	})

	t.Run("map present", func(t *testing.T) {
		m := Metadata{}
		m["test"] = "test2"
		writer := MarshalMetadata(m)
		require.NotNil(t, writer)
		b := bytes.NewBufferString("")
		writer.MarshalGQL(b)
		assert.Equal(t, "{\"test\":\"test2\"}\n", b.String())
	})
}

// TestMetadata_MarshalBSONValue will test the method MarshalBSONValue()
func TestMetadata_MarshalBSONValue(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		m := Metadata{}
		outType, outBytes, err := m.MarshalBSONValue()
		require.Equal(t, bsontype.Null, outType)
		assert.Nil(t, outBytes)
		require.NoError(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		m := Metadata{
			"test-key": "test-value",
		}
		outType, outBytes, err := m.MarshalBSONValue()
		require.NoError(t, err)
		assert.Equal(t, bsontype.Array, outType)
		assert.NotNil(t, outBytes)
		outHex := hex.EncodeToString(outBytes[:])

		out := new(map[string]interface{})
		outBytes, hexErr := hex.DecodeString(outHex)
		require.NoError(t, hexErr)
		err = bson.Unmarshal(outBytes, out)
		require.NoError(t, err)
		jsonOut, jsonErr := json.Marshal(out)
		require.NoError(t, jsonErr)
		assert.Equal(t, "{\"0\":{\"k\":\"test-key\",\"v\":\"test-value\"}}", string(jsonOut))

		// check that it is not normal marshaling
		_, inHex, _ := bson.MarshalValue(m)
		assert.NotEqual(t, hex.EncodeToString(inHex[:]), outHex)
	})
}

// TestMetadata_UnmarshalBSONValue will test the method UnmarshalBSONValue()
func TestMetadata_UnmarshalBSONValue(t *testing.T) {
	t.Parallel()

	t.Run("nil string", func(t *testing.T) {
		var m Metadata
		err := m.UnmarshalBSONValue(bsontype.Null, nil)
		require.NoError(t, err)
		assert.Nil(t, m)
	})

	t.Run("string", func(t *testing.T) {
		var m Metadata
		// this hex is a bson array [{k: "test-key", v: "test-value"}]
		b, _ := hex.DecodeString("2f000000033000270000000276000b000000746573742d76616c756500026b0009000000746573742d6b6579000000")
		err := m.UnmarshalBSONValue(bsontype.Array, b)
		require.NoError(t, err)
		assert.Equal(t, Metadata{"test-key": "test-value"}, m)
	})
}

// TestMetadata_GormDataType will test the method GormDataType()
func TestMetadata_GormDataType(t *testing.T) {
	t.Parallel()

	m := new(Metadata)
	assert.Equal(t, gormTypeText, m.GormDataType())
}

// TestMetadata_GormDBDataType will test the method GormDBDataType()
func TestMetadata_GormDBDataType(t *testing.T) {
	t.Parallel()

	t.Run("panic, no db", func(t *testing.T) {
		require.Panics(t, func() {
			m := new(Metadata)
			assert.Equal(t, datastore.JSON, m.GormDBDataType(nil, nil))
		})
	})

	t.Run("generic dialector", func(t *testing.T) {
		m := new(Metadata)
		db, err := gorm.Open(sqlite.Open(""), &gorm.Config{})
		require.NotNil(t, db)
		require.NoError(t, err)
		assert.Equal(t, datastore.JSON, m.GormDBDataType(db, nil))
	})

	t.Run("postgres dialector", func(t *testing.T) {
		/*dsn := "host=localhost user=postgres password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		require.NotNil(t, db)
		require.NoError(t, err)*/
	})
}
