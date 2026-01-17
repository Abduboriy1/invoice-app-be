// internal/domain/timeentry/entity.go
package timeentry

import (
	"time"

	"github.com/google/uuid"
)

type TimeEntry struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	InvoiceID     *uuid.UUID
	Description   string
	Hours         float64
	HourlyRate    *float64
	Date          time.Time
	JiraIssueKey  *string
	JiraWorklogID *string
	JiraSyncedAt  *time.Time
	IsBillable    bool
	IsInvoiced    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
