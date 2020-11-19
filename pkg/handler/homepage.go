package handler

import (
	"html/template"
	"log"
	"net/http"
	"saverbate/pkg/broadcast"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/authboss/v3"
)

// HomepageHandler is common handler for all actions
type HomepageHandler struct {
	db *sqlx.DB
}

// NewHomepageHandler returns ApplicationHandler interactor
func NewHomepageHandler(db *sqlx.DB) *HomepageHandler {
	return &HomepageHandler{db: db}
}

func (h *HomepageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data authboss.HTMLData
	dataIntf := r.Context().Value(authboss.CTXKeyData)

	if dataIntf == nil {
		data = authboss.HTMLData{}
	} else {
		data = dataIntf.(authboss.HTMLData)
	}

	records, err := broadcast.FeaturedRecords(h.db)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	data.MergeKV("records", records)

	template.Must(
		template.New("homepage").ParseFiles(
			"web/templates/layout.html",
			"web/templates/homepage.html",
		),
	).ExecuteTemplate(w, "layout.html", data)
}
