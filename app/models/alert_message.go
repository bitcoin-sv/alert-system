package models

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/bitcoin-sv/alert-system/utils"
	"github.com/bitcoinschema/go-bitcoin"
	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvutil"
	"github.com/libsv/go-bt/v2/chainhash"
	"github.com/mrz1836/go-datastore"
)

// AlertMessage is an object representing an alert message
type AlertMessage struct {
	// Base model
	model.Model `bson:",inline"`

	// Model specific fields
	ID             uint64 `json:"id" toml:"id" yaml:"id" bson:"_id" gorm:"primaryKey;comment:This is a unique identifier"`
	Hash           string `json:"hash" toml:"hash" yaml:"hash" bson:"hash" gorm:"<-;type:char(64);index;comment:This is the hash"`
	SequenceNumber uint32 `json:"sequence_number" toml:"sequence_number" yaml:"sequence_number" bson:"sequence_number" gorm:"<-;type:int8;index;comment:This is the alert sequence number"`
	Raw            string `json:"raw" toml:"raw" yaml:"raw" bson:"raw" gorm:"<-;type:text;comment:This is the raw alert message"`
	Processed      bool   `json:"processed" toml:"processed" yaml:"processed" bson:"processed" gorm:"<-;type:boolean;comment:This determine if the alert was processed"`

	// Private fields (never to be exported)
	alertType  AlertType
	data       []byte
	message    []byte
	signatures [][]byte
	timestamp  uint64
	version    uint32
}

// AlertMessageInterface is the interface for alert messages
type AlertMessageInterface interface {
	Read(msg []byte) error
	Do(ctx context.Context) error
}

// NewAlertMessage creates a new alert message
func NewAlertMessage(opts ...model.Options) *AlertMessage {
	return &AlertMessage{
		Model: *model.NewBaseModel(model.NameAlertMessage, opts...),
	}
}

// Name will get the name of the model
func (m *AlertMessage) Name() string {
	return model.NameAlertMessage.String()
}

// GetTableName will get the database table name of the model
func (m *AlertMessage) GetTableName() string {
	return model.TableAlertMessages
}

// GetID will get the model ID
func (m *AlertMessage) GetID() uint64 {
	return m.ID
}

// Display filter the model for display
func (m *AlertMessage) Display() interface{} {
	return m
}

// Migrate will run model specific migrations on startup
func (m *AlertMessage) Migrate(client datastore.ClientInterface) error {
	return client.IndexMetadata(client.GetTableName(model.TableAlertMessages), model.MetadataField)
}

// BeginSaveWithTx will start saving the model into the Datastore with the provided transaction
func (m *AlertMessage) BeginSaveWithTx(ctx context.Context, tx *datastore.Transaction) ([]model.BaseInterface, error) {
	return model.BeginSaveWithTx(ctx, tx, m)
}

// Save will save the model into the Datastore
func (m *AlertMessage) Save(ctx context.Context) error {
	return model.Save(ctx, m)
}

// SetAlertType will set the alert type
func (m *AlertMessage) SetAlertType(t AlertType) {
	m.alertType = t
}

// GetAlertType will get the alert type
func (m *AlertMessage) GetAlertType() AlertType {
	return m.alertType
}

// SetRawMessage will set the alert raw message
func (m *AlertMessage) SetRawMessage(msg []byte) {
	m.message = msg
}

// GetRawMessage will get the raw message
func (m *AlertMessage) GetRawMessage() []byte {
	return m.message
}

// GetRawData will get the raw data
func (m *AlertMessage) GetRawData() []byte {
	return m.data
}

// SerializeData serializes the data
func (m *AlertMessage) SerializeData() {
	var ret []byte
	ret = binary.LittleEndian.AppendUint32(ret, m.version)
	ret = binary.LittleEndian.AppendUint32(ret, m.SequenceNumber)
	ret = binary.LittleEndian.AppendUint64(ret, m.timestamp)
	ret = binary.LittleEndian.AppendUint32(ret, uint32(m.alertType))
	ret = append(ret, m.message...)
	m.data = ret
	m.Hash = chainhash.DoubleHashH(m.data).String()
}

// Serialize serializes the alert
func (m *AlertMessage) Serialize() []byte {
	m.SerializeData()
	data := m.data
	for _, sig := range m.signatures {
		data = append(data, sig...)
	}
	m.Raw = hex.EncodeToString(data)
	return data
}

// SetSignatures sets the signatures on the alert
func (m *AlertMessage) SetSignatures(sigs [][]byte) {
	m.signatures = sigs
}

// AreSignaturesValid checks if the signatures are valid
func (m *AlertMessage) AreSignaturesValid(ctx context.Context) (bool, error) {
	keys, err := GetActivePublicKey(ctx, nil, model.WithAllDependencies(m.Config()))
	if err != nil {
		return false, err
	} else if len(keys) == 0 {
		return false, fmt.Errorf("no active public keys found")
	}

	// Loop through all signatures
	for _, sig := range m.signatures {
		b64Sig := base64.StdEncoding.EncodeToString(sig)
		valid := false

		// Loop through all keys
		for _, key := range keys {

			// Get the public key
			var pub *bsvec.PublicKey
			if pub, err = bitcoin.PubKeyFromString(key.Key); err != nil {
				return false, err
			}

			// Get the address
			var addr *bsvutil.LegacyAddressPubKeyHash
			if addr, err = bitcoin.GetAddressFromPubKey(pub, true); err != nil {
				return false, err
			} else if addr == nil {
				return false, errors.New("failed to convert pub key to address")
			}

			// Verify the message
			if err = bitcoin.VerifyMessage(addr.String(), b64Sig, hex.EncodeToString(m.data)); err != nil {
				m.Config().Services.Log.Debugf("error verifying signature %x: %v", sig, err)
				continue
			}
			valid = true
			break
		}
		if !valid {
			return false, nil
		}
	}

	return true, nil
}

// ProcessAlertMessage processes the alert message and converts to an alert message interface
func (m *AlertMessage) ProcessAlertMessage() AlertMessageInterface {
	switch m.alertType {
	case AlertTypeInformational:
		return &AlertMessageInformational{
			AlertMessage: *m,
		}
	case AlertTypeFreezeUtxo:
		return &AlertMessageFreezeUtxo{
			AlertMessage: *m,
		}
	case AlertTypeUnfreezeUtxo:
		return &AlertMessageUnfreezeUtxo{
			AlertMessage: *m,
		}
	case AlertTypeConfiscateUtxo:
		return &AlertMessageConfiscateTransaction{
			AlertMessage: *m,
		}
	case AlertTypeBanPeer:
		return &AlertMessageBanPeer{
			AlertMessage: *m,
		}
	case AlertTypeUnbanPeer:
		return &AlertMessageUnbanPeer{
			AlertMessage: *m,
		}
	case AlertTypeInvalidateBlock:
		return &AlertMessageInvalidateBlock{
			AlertMessage: *m,
		}
	case AlertTypeSetKeys:
		return &AlertMessageSetKeys{
			AlertMessage: *m,
			Hash:         m.Hash,
		}
	default:
		return nil
	}
}

// SetVersion sets the version of the message
func (m *AlertMessage) SetVersion(ver uint32) {
	m.version = ver
}

// Version returns the version of the message
func (m *AlertMessage) Version() uint32 {
	return m.version
}

// SetTimestamp sets the timestamp of the message
func (m *AlertMessage) SetTimestamp(ts uint64) {
	m.timestamp = ts
}

// Timestamp returns the timestamp of the message
func (m *AlertMessage) Timestamp() uint64 {
	return m.timestamp
}

// ReadRaw sets the model fields based on the raw message
func (m *AlertMessage) ReadRaw() error {
	if len(m.GetRawMessage()) == 0 {
		ak, err := hex.DecodeString(m.Raw)
		if err != nil {
			return err
		}
		m.SetRawMessage(ak)
	}

	if len(m.GetRawMessage()) < 16 {
		// todo DETERMINE ACTUAL PROPER LENGTH
		return fmt.Errorf("alert needs to be at least 16 bytes")
	}
	ak := m.GetRawMessage()
	version := binary.LittleEndian.Uint32(ak[:4])
	sequenceNumber := binary.LittleEndian.Uint32(ak[4:8])
	timestamp := binary.LittleEndian.Uint64(ak[8:16])
	alertType := binary.LittleEndian.Uint32(ak[16:20])

	alertAndSignature := ak[20:]

	// Assume 3 signatures, maybe disable alert will require 2 (0x09)
	sigLen := 195
	switch alertType {
	case uint32(99):
		sigLen = 128
	}

	// This is the minimum length this data should be. Signature byte length + 2 bytes
	// This would imply an informational alert with a message 1 byte long... not practical
	// but possible. Regardless let's just error out now if this length is lower. At least
	// allows us to grab the expected signature.
	if len(alertAndSignature) < sigLen+2 {
		return fmt.Errorf("alert message is invalid - too short length")
	}

	// Get alert message bytes
	alert := alertAndSignature[:len(alertAndSignature)-sigLen]

	// Get signature bytes
	signatures := alertAndSignature[len(alertAndSignature)-sigLen:]
	var sigs [][]byte

	// Loop through all signatures and create array
	for i := 0; i < sigLen/65; i++ {
		sigs = append(sigs, signatures[:65])
		signatures = signatures[65:]
	}

	dataLen := 20 + len(alert)

	m.SetAlertType(AlertType(alertType))
	m.message = alert
	m.SequenceNumber = sequenceNumber
	m.timestamp = timestamp
	m.version = version
	m.data = ak[:dataLen]
	m.signatures = sigs
	_ = m.Serialize()
	return nil
}

// NewAlertFromBytes creates a new alert from bytes
func NewAlertFromBytes(ak []byte, opts ...model.Options) (*AlertMessage, error) {
	opts = append(opts, model.New())
	newAlert := NewAlertMessage(opts...)
	newAlert.SetRawMessage(ak)
	err := newAlert.ReadRaw()
	if err != nil {
		return nil, err
	}

	// Return alert
	return newAlert, nil
}

// GetAlertMessageBySequenceNumber will get the model with the given conditions
func GetAlertMessageBySequenceNumber(ctx context.Context, sequenceNumber uint32, opts ...model.Options) (*AlertMessage, error) {

	// Get the record
	message := NewAlertMessage(opts...)
	message.SequenceNumber = sequenceNumber
	conditions := make(map[string]interface{})
	conditions["sequence_number"] = sequenceNumber
	if err := model.Get(
		ctx, message, conditions, model.DefaultDatabaseReadTimeout, true, // In-case an update is occurring
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return message, nil
}

// GetLatestAlert will get the model with the given conditions
func GetLatestAlert(ctx context.Context, metadata *model.Metadata, opts ...model.Options) (*AlertMessage, error) {

	// Set the conditions
	conditions := &map[string]interface{}{
		utils.FieldDeletedAt: map[string]interface{}{ // IS NULL
			utils.ExistsCondition: false,
		},
	}

	// Set the query params
	queryParams := &datastore.QueryParams{
		Page:          1,
		PageSize:      1,
		OrderByField:  utils.FieldSequenceNumber,
		SortDirection: utils.SortDescending,
	}

	// Get the record
	modelItems := make([]*AlertMessage, 0)
	if err := model.GetModelsByConditions(
		ctx, model.NameAlertMessage, &modelItems, metadata, conditions, queryParams, opts...,
	); err != nil {
		return nil, err
	} else if len(modelItems) == 0 {
		return nil, nil
	}

	// Return the first item (only item)
	return modelItems[0], nil
}

// GetAllAlerts
func GetAllAlerts(ctx context.Context, metadata *model.Metadata, opts ...model.Options) ([]*AlertMessage, error) {
	// Set the conditions
	conditions := &map[string]interface{}{
		utils.FieldDeletedAt: map[string]interface{}{ // IS NULL
			utils.ExistsCondition: false,
		},
	}

	// Set the query params
	queryParams := &datastore.QueryParams{
		OrderByField:  utils.FieldSequenceNumber,
		SortDirection: utils.SortAscending,
	}

	// Get the record
	modelItems := make([]*AlertMessage, 0)
	if err := model.GetModelsByConditions(
		ctx, model.NameAlertMessage, &modelItems, metadata, conditions, queryParams, opts...,
	); err != nil {
		return nil, err
	} else if len(modelItems) == 0 {
		return nil, nil
	}

	// Return the first item (only item)
	return modelItems, nil
}

// GetAllUnprocessedAlerts will get all alerts that weren't successfully processed
func GetAllUnprocessedAlerts(ctx context.Context, metadata *model.Metadata, opts ...model.Options) ([]*AlertMessage, error) {

	// Set the conditions
	conditions := &map[string]interface{}{
		utils.FieldDeletedAt: map[string]interface{}{ // IS NULL
			utils.ExistsCondition: false,
		},
		"processed": false,
	}

	// Set the query params
	queryParams := &datastore.QueryParams{
		OrderByField:  utils.FieldSequenceNumber,
		SortDirection: utils.SortAscending,
	}

	// Get the record
	modelItems := make([]*AlertMessage, 0)
	if err := model.GetModelsByConditions(
		ctx, model.NameAlertMessage, &modelItems, metadata, conditions, queryParams, opts...,
	); err != nil {
		return nil, err
	} else if len(modelItems) == 0 {
		return nil, nil
	}

	// Return the first item (only item)
	return modelItems, nil
}
