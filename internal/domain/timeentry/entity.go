// internal/domain/timeentry/entity.go
package timeentry

import (
	"time"

	"github.com/google/uuid"
)

type TimeEntry struct {
	ID            uuid.UUID  `db:"id"`
	UserID        uuid.UUID  `db:"user_id"`
	InvoiceID     *uuid.UUID `db:"invoice_id"`
	Description   string     `db:"description"`
	Hours         float64    `db:"hours"`
	HourlyRate    *float64   `db:"hourly_rate"`
	Date          time.Time  `db:"date"`
	JiraIssueKey  *string    `db:"jira_issue_key"`
	JiraWorklogID *string    `db:"jira_worklog_id"`
	JiraSyncedAt  *time.Time `db:"jira_synced_at"`
	IsBillable    bool       `db:"is_billable"`
	IsInvoiced    bool       `db:"is_invoiced"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}
