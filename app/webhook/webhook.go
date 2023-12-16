// Package webhook provides a webhook client for sending alerts to a webhook URL
package webhook

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/tokenized/pkg/json"
)

// Payload is the payload for the webhook
type Payload struct {
	AlertType models.AlertType `json:"alert_type"`
	Raw       string           `json:"raw"`
	Sequence  uint32           `json:"sequence"`
	Text      string           `json:"text"`
}

// PostAlert sends an alert to a webhook URL using the provided http client
func PostAlert(ctx context.Context, httpClient config.HTTPInterface, url string, alert *models.AlertMessage) error {

	// Validate the URL length
	if len(url) == 0 {
		return fmt.Errorf("webhook URL is not configured")
	}

	// Validate the URL prefix
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("webhook URL [%s] is does not have a valid prefix", url)
	}

	// Serialize the alert
	raw := alert.Serialize()

	// Create the payload
	p := Payload{
		AlertType: alert.GetAlertType(),
		Sequence:  alert.SequenceNumber,
		Raw:       hex.EncodeToString(raw),
		Text:      fmt.Sprintf("received alert type [%d], sequence [%d], with raw data [%x]", alert.GetAlertType(), alert.SequenceNumber, raw),
	}

	// Marshal the payload
	var payload []byte
	var err error
	if payload, err = json.Marshal(p); err != nil {
		return err
	}

	// Create the http request
	var req *http.Request
	if req, err = http.NewRequestWithContext(
		ctx, http.MethodPost, url, bytes.NewReader(payload),
	); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Fire the http request
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		return err
	}
	defer func() {
		if res != nil && res.Body != nil {
			_ = res.Body.Close()
		}
	}()

	// Validate the response
	if res != nil && res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code [%d] sending payload to webhook", res.StatusCode)
	}
	return nil
}
