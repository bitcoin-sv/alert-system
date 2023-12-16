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

// HealthResponse is the response for the health endpoint
type HealthResponse struct {
	Alert    models.AlertMessage `json:"alert"`
	Sequence uint32              `json:"sequence"`
	Synced   bool                `json:"synced"`
}

// health will return the health of the API and the current alert
func (a *Action) health(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the latest alert
	alert, err := models.GetLatestAlert(req.Context(), nil, model.WithAllDependencies(a.Config))
	if err != nil {
		app.APIErrorResponse(w, req, http.StatusBadRequest, err)
		return
	} else if alert == nil {
		app.APIErrorResponse(w, req, http.StatusNotFound, errors.New("alert not found"))
		return
	}

	// Return the response
	_ = apirouter.ReturnJSONEncode(
		w,
		http.StatusOK,
		json.NewEncoder(w),
		HealthResponse{
			Alert:    *alert,
			Sequence: alert.SequenceNumber,
			Synced:   true, // TODO actually fetch this state from the DB somehow, or from the server struct
		}, []string{"alert", "synced", "sequence"})
}
