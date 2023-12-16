package base

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bitcoin-sv/alert-system/app/webhook"

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
		app.APIErrorResponse(w, req, http.StatusInternalServerError, err)
		return
	} else if alertModel == nil {
		app.APIErrorResponse(w, req, http.StatusNotFound, errors.New("alert not found"))
		return
	}
	err = alertModel.ReadRaw()
	if err != nil {
		app.APIErrorResponse(w, req, http.StatusInternalServerError, errors.New("alert faile"))
		return
	}
	am := alertModel.ProcessAlertMessage()
	if am == nil {
		app.APIErrorResponse(w, req, http.StatusInternalServerError, errors.New("alert not valid type"))
		return
	}
	err = am.Read(alertModel.GetRawMessage())
	if err != nil {
		app.APIErrorResponse(w, req, http.StatusInternalServerError, err)
		return
	}
	p := webhook.Payload{
		AlertType: alertModel.GetAlertType(),
		Sequence:  alertModel.SequenceNumber,
		Raw:       hex.EncodeToString(alertModel.GetRawData()),
		Text:      am.MessageString(),
	}
	// Return the response
	_ = apirouter.ReturnJSONEncode(
		w,
		http.StatusOK,
		json.NewEncoder(w),
		p, []string{"sequence", "raw", "text", "alert_type"})
}
