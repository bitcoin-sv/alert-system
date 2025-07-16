package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bsv-blockchain/go-sdk/util"
)

// AlertMessageInformational is an informational alert
type AlertMessageInformational struct {
	AlertMessage
	MessageLength uint64 `json:"message_length"`
	Message       []byte `json:"message"`
}

// Read reads the alert message from the byte slice
func (a *AlertMessageInformational) Read(alert []byte) error {
	reader := util.NewReader(alert[:])

	// read the message length
	length, err := reader.ReadVarInt()
	if err != nil {
		return err
	}
	if length > uint64(len(reader.Data)) {
		return errors.New("info message length is longer than buffer")
	}

	// read the message
	var msg []byte
	for i := uint64(0); i < length; i++ {
		var b byte
		if b, err = reader.ReadByte(); err != nil {
			return fmt.Errorf("failed to read message: %s", err.Error())
		}
		msg = append(msg, b)
	}
	if !reader.IsComplete() {
		return fmt.Errorf("too many bytes in alert message")
	}
	a.Message = msg
	a.MessageLength = length
	return nil
}

// Do execute the alert
func (a *AlertMessageInformational) Do(_ context.Context) error {
	a.Config().Services.Log.Infof("[informational alert]: %s", a.Message)
	return nil
}

// ToJSON is the alert in JSON format
func (a *AlertMessageInformational) ToJSON(_ context.Context) []byte {
	m := a.ProcessAlertMessage()
	// TODO: Come back and add a message interface for each alert
	_ = m.Read(a.GetRawMessage())
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return []byte{}
	}
	return data
}

// MessageString executes the alert
func (a *AlertMessageInformational) MessageString() string {
	return fmt.Sprintf("Informational: %s", a.Message)
}
