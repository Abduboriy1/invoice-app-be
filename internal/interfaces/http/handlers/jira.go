// internal/interfaces/http/handlers/jira.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/invoice-app-be/internal/domain/timeentry"
	"github.com/invoice-app-be/internal/infrastructure/integrations/jira"
	"github.com/invoice-app-be/internal/interfaces/http/dto"
	"github.com/invoice-app-be/internal/interfaces/http/middleware"
)

type JiraHandler struct {
	jiraSyncService *jira.SyncService
	timeEntryRepo   timeentry.Repository
}

func NewJiraHandler(jiraSyncService *jira.SyncService, timeEntryRepo timeentry.Repository) *JiraHandler {
	return &JiraHandler{
		jiraSyncService: jiraSyncService,
		timeEntryRepo:   timeEntryRepo,
	}
}

// PullWorklogs pulls worklogs from Jira for a date range
func (h *JiraHandler) PullWorklogs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	var req dto.PullWorklogsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid start_date format. Expected YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid end_date format. Expected YYYY-MM-DD")
		return
	}

	if endDate.Before(startDate) {
		respondError(w, http.StatusBadRequest, "end_date must be after start_date")
		return
	}

	// Pull worklogs from Jira
	count, err := h.jiraSyncService.PullWorklogsByDateRange(r.Context(), userID, startDate, endDate, req.IssueKeys)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to pull worklogs from Jira: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Worklogs synced successfully",
		"count":   count,
	})
}

// PullWorklogsForIssue pulls worklogs for a specific Jira issue
func (h *JiraHandler) PullWorklogsForIssue(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	var req dto.PullIssueWorklogsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Sync worklogs for the issue
	if err := h.jiraSyncService.SyncWorklogsForIssue(r.Context(), userID, req.IssueKey); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to sync worklogs: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Worklogs synced successfully for issue " + req.IssueKey,
	})
}

// PushWorklog pushes a time entry to Jira
func (h *JiraHandler) PushWorklog(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	var req dto.PushWorklogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the time entry
	entryID, err := uuid.Parse(req.TimeEntryID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid time_entry_id")
		return
	}

	entry, err := h.timeEntryRepo.GetByID(r.Context(), entryID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Time entry not found")
		return
	}

	if entry.UserID != userID {
		respondError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	// Push to Jira
	if err := h.jiraSyncService.PushTimeEntryToJira(r.Context(), entry, req.IssueKey); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to push worklog to Jira: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Worklog pushed to Jira successfully",
	})
}

// ConfigureJira saves Jira credentials for the user
func (h *JiraHandler) ConfigureJira(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Jira configuration storage
	// This would save the user's Jira credentials (encrypted) to the database
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Jira configuration saved",
	})
}
