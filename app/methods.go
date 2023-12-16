package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// Head is a basic response for any generic HEAD request
func Head(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

// NotFound handles all 404 requests
func NotFound(w http.ResponseWriter, req *http.Request) {
	apirouter.ReturnResponse(w, req, http.StatusNotFound, "route not found: "+req.RequestURI)
}

// MethodNotAllowed handles all 405 requests
func MethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	apirouter.ReturnResponse(w, req, http.StatusMethodNotAllowed, req.RequestURI+":"+req.Method)
}
