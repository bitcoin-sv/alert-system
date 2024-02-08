package base

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bitcoin-sv/alert-system/app"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/bitcoin-sv/alert-system/app/models/model"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// alerts will return the saved
func (a *Action) alert(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Read params
	params := apirouter.GetParams(req)
	if params == nil {
		apiError := apirouter.ErrorFromRequest(req, "parameters is nil", "no parameters specified", http.StatusBadRequest, http.StatusBadRequest, "")
		apirouter.ReturnResponse(w, req, apiError.Code, apiError)
		return
	}
	idStr := params.GetString("sequence")
	if idStr == "" {
		apiError := apirouter.ErrorFromRequest(req, "missing sequence param", "missing sequence param", http.StatusBadRequest, http.StatusBadRequest, "")
		apirouter.ReturnResponse(w, req, apiError.Code, apiError)
		return
	}
	sequenceNumber, err := strconv.Atoi(idStr)
	if err != nil {
		apiError := apirouter.ErrorFromRequest(req, "sequence is invalid", "sequence is invalid", http.StatusBadRequest, http.StatusBadRequest, "")
		apirouter.ReturnResponse(w, req, apiError.Code, apiError)
		return
	}

	// Get alert
	alertModel, err := models.GetAlertMessageBySequenceNumber(req.Context(), uint32(sequenceNumber), model.WithAllDependencies(a.Config))
	if err != nil {
		app.APIErrorResponse(w, req, http.StatusBadRequest, err)
		return
	} else if alertModel == nil {
		app.APIErrorResponse(w, req, http.StatusNotFound, errors.New("alert not found"))
		return
	}

	// Return the response
	_ = apirouter.ReturnJSONEncode(
		w,
		http.StatusOK,
		json.NewEncoder(w),
		alertModel, []string{"sequence_number", "raw"})
}
