// internal/infrastructure/integrations/jira/mapper.go
package jira

import (
	"time"

	"github.com/google/uuid"
	"github.com/invoice-app-be/internal/domain/timeentry"
)

// MapWorklogToTimeEntry converts a Jira worklog to a domain TimeEntry
func MapWorklogToTimeEntry(userID uuid.UUID, worklog Worklog) *timeentry.TimeEntry {
	// Convert time spent seconds to hours
	hours := float64(worklog.TimeSpentSeconds) / 3600.0

	// Use the Started field directly (it's already time.Time)
	date := worklog.Started.Truncate(24 * time.Hour)

	// Store Jira-specific fields
	issueKey := worklog.IssueKey
	worklogID := worklog.ID
	syncedAt := time.Now()

	return &timeentry.TimeEntry{
		ID:            uuid.New(),
		UserID:        userID,
		InvoiceID:     nil,
		Description:   worklog.Comment,
		Hours:         hours,
		HourlyRate:    nil,
		Date:          date,
		JiraIssueKey:  &issueKey,
		JiraWorklogID: &worklogID,
		JiraSyncedAt:  &syncedAt,
		IsBillable:    true,
		IsInvoiced:    false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// MapTimeEntryToJiraWorklog converts a TimeEntry to Jira worklog format
func MapTimeEntryToJiraWorklog(entry *timeentry.TimeEntry) (issueKey string, timeSpentSeconds int, comment string, started time.Time) {
	// Convert hours to seconds
	timeSpentSeconds = int(entry.Hours * 3600)

	// Use the entry's date as the start time (at midnight)
	started = entry.Date

	// Use description as comment
	comment = entry.Description

	// Get issue key if available
	if entry.JiraIssueKey != nil {
		issueKey = *entry.JiraIssueKey
	}

	return issueKey, timeSpentSeconds, comment, started
}
