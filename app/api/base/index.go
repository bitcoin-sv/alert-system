package base

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// index is the default index route of the API for testing purposes: (Hello World)
func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	apirouter.ReturnResponse(
		w, req, http.StatusOK, "Bitcoin SV Alert System",
	)
}
