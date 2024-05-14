package base

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// PeersResponse is the response for the health endpoint
type PeersResponse struct {
	Peers []string `json:"peers"`
}

// health will return the health of the API and the current alert
func (a *Action) peers(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	peers := a.P2pServer.Peers()

	// Return the response
	_ = apirouter.ReturnJSONEncode(
		w,
		http.StatusOK,
		json.NewEncoder(w),
		PeersResponse{
			Peers: peers,
		}, []string{"peers"})
}
