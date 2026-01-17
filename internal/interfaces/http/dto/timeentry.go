package dto

import "time"

type CreateTimeEntryRequest struct {
	ProjectID   string    `json:"project_id"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Billable    bool      `json:"billable"`
}

type UpdateTimeEntryRequest struct {
	ProjectID   string    `json:"project_id"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Billable    bool      `json:"billable"`
}

type TimeEntryResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProjectID   string    `json:"project_id"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    int64     `json:"duration_seconds"`
	Billable    bool      `json:"billable"`
	InvoiceID   *string   `json:"invoice_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
