package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"saverbate/pkg/broadcast"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/authboss/v3"
)

const perPage = 12

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

	nextPaginationState := "#"
	prevPaginationState := "#"
	page := currentPage(r)
	total, err := broadcast.TotalFeaturedRecords(h.db)
	if err != nil && err != sql.ErrNoRows {
		log.Panicf("ERROR: %v", err)
	}

	records, err := broadcast.FeaturedRecords(h.db, page, perPage)
	if err != nil && err != sql.ErrNoRows {
		log.Panicf("ERROR: %v", err)
	}

	totalPages := int64(1)
	if total > perPage {
		totalPages = int64(math.Ceil(float64(total) / float64(perPage)))
	}

	if page > 1 {
		prevPaginationState = fmt.Sprintf("/?page=%d", page-1)
	}
	if page < totalPages {
		nextPaginationState = fmt.Sprintf("/?page=%d", page+1)
	}

	data.MergeKV("records", records)
	data.MergeKV("prevPaginationState", prevPaginationState)
	data.MergeKV("nextPaginationState", nextPaginationState)

	template.Must(
		template.New("homepage").ParseFiles(
			"web/templates/layout.html",
			"web/templates/homepage.html",
		),
	).ExecuteTemplate(w, "layout.html", data)
}

func currentPage(r *http.Request) int64 {
	p := fetchParamValue(r, "page")
	if p == "" {
		return 1
	}

	page, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return 1
	}

	return page
}

func fetchParamValue(r *http.Request, key string) string {
	val, ok := r.URL.Query()[key]
	if !ok || len(val[0]) == 0 {
		return ""
	}

	return val[0]
}
