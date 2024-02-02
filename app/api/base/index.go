package base

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/bitcoin-sv/alert-system/app/models/model"

	"github.com/bitcoin-sv/alert-system/app/models"

	"github.com/julienschmidt/httprouter"
)

//go:embed ui/templates/*
var content embed.FS

// PageData contains the page data
type PageData struct {
	Alerts []*models.AlertMessage
}

func substr(s string, start, length int) string {
	end := start + length
	if start < 0 || start >= len(s) || end > len(s) {
		return s
	}
	return s[start:end]
}

// index is the default index route of the API for testing purposes: (Hello World)
func (a *Action) index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	htmlContent, err := content.ReadFile("ui/templates/index.tmpl")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ts, err := template.New("index").Funcs(template.FuncMap{"substr": substr}).Parse(string(htmlContent))
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	alerts, err := models.GetAllAlerts(context.Background(), nil, model.WithAllDependencies(a.Config))
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := PageData{
		Alerts: alerts,
	}

	// Then we use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.
	err = ts.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
