package base

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bitcoin-sv/alert-system/app"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// AlertsResponse is the response for the alerts endpoint
type AlertsResponse struct {
	Alerts         []*models.AlertMessage `json:"alerts"`
	LatestSequence uint32                 `json:"latest_sequence"`
}

// alerts will return the saved
func (a *Action) alerts(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get all alerts
	alerts, err := models.GetAllAlerts(req.Context(), nil, model.WithAllDependencies(a.Config))
	if err != nil {
		app.APIErrorResponse(w, req, http.StatusBadRequest, err)
		return
	} else if alerts == nil {
		app.APIErrorResponse(w, req, http.StatusNotFound, errors.New("alert not found"))
		return
	}

	// Return the response
	_ = apirouter.ReturnJSONEncode(
		w,
		http.StatusOK,
		json.NewEncoder(w),
		AlertsResponse{
			Alerts:         alerts,
			LatestSequence: alerts[len(alerts)-1].SequenceNumber,
		}, []string{"alerts", "latest_sequence"})
}
