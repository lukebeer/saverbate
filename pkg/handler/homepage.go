package handler

import (
	"html/template"
	"net/http"

	"github.com/volatiletech/authboss/v3"
)

// HomepageHandler is common handler for all actions
type HomepageHandler struct {
}

// NewHomepageHandler returns ApplicationHandler interactor
func NewHomepageHandler() *HomepageHandler {
	return &HomepageHandler{}
}

func (h *HomepageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data authboss.HTMLData
	dataIntf := r.Context().Value(authboss.CTXKeyData)

	if dataIntf == nil {
		data = authboss.HTMLData{}
	} else {
		data = dataIntf.(authboss.HTMLData)
	}

	template.Must(
		template.New("homepage").ParseFiles(
			"web/templates/layout.html",
			"web/templates/homepage.html",
		),
	).ExecuteTemplate(w, "layout.html", data)
}
