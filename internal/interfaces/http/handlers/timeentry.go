// internal/interfaces/http/handlers/timeentry.go
package handlers

import (
	"net/http"

	"github.com/invoice-app-be/internal/domain/timeentry"
)

type TimeEntryHandler struct {
	service *timeentry.Service
}

func NewTimeEntryHandler(service *timeentry.Service) *TimeEntryHandler {
	return &TimeEntryHandler{service: service}
}

func (h *TimeEntryHandler) List(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, []string{})
}

func (h *TimeEntryHandler) Create(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *TimeEntryHandler) Get(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *TimeEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *TimeEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusNoContent, nil)
}

func (h *TimeEntryHandler) SyncToJira(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "synced"})
}
