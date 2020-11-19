package handler

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"saverbate/pkg/broadcast"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/authboss/v3"
)

// ShowHandler shows video
type ShowHandler struct {
	db *sqlx.DB
}

// NewShowHandler returns ShowHandler interactor
func NewShowHandler(db *sqlx.DB) *ShowHandler {
	return &ShowHandler{db: db}
}

func (h *ShowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recordUUID := chi.URLParam(r, "uuid")

	record, err := broadcast.FindByUUID(h.db, recordUUID)
	if err == sql.ErrNoRows {
		NotFoundError(w, r)
		return
	} else if err != nil {
		log.Panicf("ERROR: %v", err)
	}

	var data authboss.HTMLData
	dataIntf := r.Context().Value(authboss.CTXKeyData)

	if dataIntf == nil {
		data = authboss.HTMLData{}
	} else {
		data = dataIntf.(authboss.HTMLData)
	}

	data.MergeKV("record", record)

	template.Must(
		template.New("homepage").ParseFiles(
			"web/templates/layout.html",
			"web/templates/show.html",
		),
	).ExecuteTemplate(w, "layout.html", data)
}
