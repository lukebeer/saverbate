package handler

import (
	"html/template"
	"log"
	"net/http"
	"saverbate/pkg/broadcast"
	"strconv"
	"time"

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

	prevRecord := recordFromParams(r, "p")
	records, err := broadcast.FeaturedRecords(h.db, prevRecord, perPage)
	if err != nil {
		log.Panicf("ERROR: %v", err)
	}

	data.MergeKV("records", records)

	// TODO: use UUID
	currentPage := currentPageParam(r)
	firstRecord := recordFromParams(r, "f")

	if currentPage > 1 && prevRecord != nil {
		//
	}

	data.MergeKV("firstRecord", prevRecord)

	// Pass last record as previous record for paginate
	l := len(records)
	if l > 0 && l == perPage {
		data.MergeKV("prevRecord", records[len(records)-1])
		data.MergeKV("nextPage", currentPage+1)
	}

	template.Must(
		template.New("homepage").ParseFiles(
			"web/templates/layout.html",
			"web/templates/homepage.html",
		),
	).ExecuteTemplate(w, "layout.html", data)
}

func recordFromParams(r *http.Request, prefix string) *broadcast.Record {
	sPrevID := fetchParamValue(r, prefix+"_i")
	if sPrevID == "" {
		return nil
	}
	prevID, err := strconv.ParseInt(sPrevID, 10, 64)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil
	}

	sPrevFollowers := fetchParamValue(r, prefix+"_f")
	if sPrevFollowers == "" {
		return nil
	}
	prevFollowers, err := strconv.ParseInt(sPrevFollowers, 10, 64)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil
	}

	sPrevFinishAt := fetchParamValue(r, prefix+"_ft")
	if sPrevFinishAt == "" {
		return nil
	}
	prevFinishAt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", sPrevFinishAt)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil
	}

	return &broadcast.Record{
		ID:        prevID,
		Followers: prevFollowers,
		FinishAt:  prevFinishAt,
	}
}

func currentPageParam(r *http.Request) int64 {
	p := fetchParamValue(r, "page")
	if p == "" {
		return 0
	}

	page, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil
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
