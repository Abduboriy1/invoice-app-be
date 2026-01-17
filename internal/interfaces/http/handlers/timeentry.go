// internal/interfaces/http/handlers/timeentry.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/invoice-app-be/internal/domain/timeentry"
	"github.com/invoice-app-be/internal/interfaces/http/dto"
	"github.com/invoice-app-be/internal/interfaces/http/middleware"
)

type TimeEntryHandler struct {
	service *timeentry.Service
}

func NewTimeEntryHandler(service *timeentry.Service) *TimeEntryHandler {
	return &TimeEntryHandler{service: service}
}

func (h *TimeEntryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	entries, err := h.service.ListTimeEntries(r.Context(), userID)
	if err != nil {
		log.Print(err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch time entries")
		return
	}

	response := make([]dto.TimeEntryResponse, len(entries))
	for i, entry := range entries {
		response[i] = dto.TimeEntryFromDomain(&entry)
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *TimeEntryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	var req dto.CreateTimeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse date - FIXED: Use date format instead of datetime
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid date format. Expected YYYY-MM-DD")
		return
	}

	// Map to domain request
	domainReq := timeentry.CreateTimeEntryRequest{
		Description: req.Description,
		Hours:       req.Hours,
		Date:        date,
		IsBillable:  req.IsBillable,
	}

	entry, err := h.service.CreateTimeEntry(r.Context(), userID, domainReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create time entry")
		return
	}

	respondJSON(w, http.StatusCreated, dto.TimeEntryFromDomain(entry))
}

func (h *TimeEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	entryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid time entry ID")
		return
	}

	var req dto.UpdateTimeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse date - FIXED: Use date format instead of datetime
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid date format. Expected YYYY-MM-DD")
		return
	}

	entry, err := h.service.UpdateTimeEntry(r.Context(), userID, entryID, timeentry.UpdateTimeEntryRequest{
		Description: req.Description,
		Hours:       req.Hours,
		Date:        date,
		IsBillable:  req.IsBillable,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update time entry")
		return
	}

	respondJSON(w, http.StatusOK, dto.TimeEntryFromDomain(entry))
}

func (h *TimeEntryHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	entryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid time entry ID")
		return
	}

	entry, err := h.service.GetTimeEntry(r.Context(), userID, entryID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Time entry not found")
		return
	}

	respondJSON(w, http.StatusOK, dto.TimeEntryFromDomain(entry))
}

func (h *TimeEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	entryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid time entry ID")
		return
	}

	if err := h.service.DeleteTimeEntry(r.Context(), userID, entryID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete time entry")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

func (h *TimeEntryHandler) SyncToJira(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	entryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid time entry ID")
		return
	}

	var req dto.SyncToJiraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.IssueKey == "" {
		respondError(w, http.StatusBadRequest, "Jira issue key is required")
		return
	}

	if err := h.service.SyncToJira(r.Context(), userID, entryID, req.IssueKey); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to sync to Jira: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "synced",
		"message": "Time entry synced to Jira successfully",
	})
}
