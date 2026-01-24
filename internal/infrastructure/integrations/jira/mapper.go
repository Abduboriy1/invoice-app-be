// internal/infrastructure/integrations/jira/mapper.go
package jira

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/invoice-app-be/internal/domain/timeentry"
)

// MapWorklogToTimeEntry converts a Jira worklog to a domain TimeEntry
func MapWorklogToTimeEntry(userID uuid.UUID, worklog Worklog) *timeentry.TimeEntry {
	// Convert time spent seconds to hours
	hours := float64(worklog.TimeSpentSeconds) / 3600.0

	// Use the Started field (it's already time.Time)
	date := worklog.Started.Truncate(24 * time.Hour)

	// Extract comment text from the nested structure
	commentText := extractCommentText(worklog.Comment)

	// Store Jira-specific fields
	issueKey := worklog.IssueKey
	worklogID := worklog.ID
	syncedAt := time.Now()

	return &timeentry.TimeEntry{
		ID:            uuid.New(),
		UserID:        userID,
		InvoiceID:     nil,
		Description:   commentText,
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

// extractCommentText extracts plain text from Jira's nested comment structure
func extractCommentText(comment Comment) string {
	var texts []string

	for _, content := range comment.Content {
		for _, textContent := range content.Content {
			if textContent.Text != "" {
				texts = append(texts, textContent.Text)
			}
		}
	}

	if len(texts) == 0 {
		return "No description"
	}

	return strings.Join(texts, " ")
}
