// internal/interfaces/http/dto/timeentry.go
package dto

import (
	"time"

	"github.com/invoice-app-be/internal/domain/timeentry"
)

type CreateTimeEntryRequest struct {
	Description string  `json:"description" validate:"required"`
	Hours       float64 `json:"hours" validate:"required,gt=0"`
	Date        string  `json:"date"`
	IsBillable  bool    `json:"is_billable"`
}

type UpdateTimeEntryRequest struct {
	Description string  `json:"description" validate:"required"`
	Hours       float64 `json:"hours" validate:"required,gt=0"`
	Date        string  `json:"date" validate:"required"`
	IsBillable  bool    `json:"is_billable"`
}

type SyncToJiraRequest struct {
	IssueKey string `json:"issue_key" validate:"required"`
}

type TimeEntryResponse struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	InvoiceID     *string  `json:"invoice_id,omitempty"`
	Description   string   `json:"description"`
	Hours         float64  `json:"hours"`
	HourlyRate    *float64 `json:"hourly_rate,omitempty"`
	Date          string   `json:"date"`
	JiraIssueKey  *string  `json:"jira_issue_key,omitempty"`
	JiraWorklogID *string  `json:"jira_worklog_id,omitempty"`
	IsBillable    bool     `json:"is_billable"`
	IsInvoiced    bool     `json:"is_invoiced"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

func TimeEntryFromDomain(entry *timeentry.TimeEntry) TimeEntryResponse {
	resp := TimeEntryResponse{
		ID:          entry.ID.String(),
		UserID:      entry.UserID.String(),
		Description: entry.Description,
		Hours:       entry.Hours,
		Date:        entry.Date.Format("2006-01-02"),
		IsBillable:  entry.IsBillable,
		IsInvoiced:  entry.IsInvoiced,
		CreatedAt:   entry.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   entry.UpdatedAt.Format(time.RFC3339),
	}

	if entry.InvoiceID != nil {
		invoiceID := entry.InvoiceID.String()
		resp.InvoiceID = &invoiceID
	}

	if entry.HourlyRate != nil {
		resp.HourlyRate = entry.HourlyRate
	}

	if entry.JiraIssueKey != nil {
		resp.JiraIssueKey = entry.JiraIssueKey
	}

	if entry.JiraWorklogID != nil {
		resp.JiraWorklogID = entry.JiraWorklogID
	}

	return resp
}
