package handler

import (
	"log"
	"net/http"

	"saverbate/pkg/utils"
)

// NotFoundHandler shows 404 page
type NotFoundHandler struct{}

// NewNotFoundHandler returns NotFoundHandler interactor
func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	NotFoundError(w, r)
}

// NotFoundError returns 404 page
func NotFoundError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	staticDir, err := utils.StaticDir()
	if err != nil {
		log.Panicf("ERROR: %v", err)
	}

	if err := utils.ServeStaticFile(staticDir+"/404.html", w); err != nil {
		log.Panicf("ERROR: %v", err)
	}
}
