// internal/infrastructure/integrations/jira/sync.go
package jira

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/invoice-app-be/internal/domain/timeentry"
)

type SyncService struct {
	client        *Client
	timeEntryRepo timeentry.Repository
}

func NewSyncService(client *Client, repo timeentry.Repository) *SyncService {
	return &SyncService{
		client:        client,
		timeEntryRepo: repo,
	}
}

func (s *SyncService) SyncWorklogsForIssue(ctx context.Context, userID uuid.UUID, issueKey string) error {
	worklogs, err := s.client.GetWorklogs(ctx, issueKey)
	if err != nil {
		return fmt.Errorf("fetching worklogs: %w", err)
	}

	for _, wl := range worklogs {
		// Check if already synced by getting all entries and checking manually
		// (since GetByJiraWorklogID doesn't exist yet)
		entries, err := s.timeEntryRepo.GetByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("checking existing entries: %w", err)
		}

		// Check if this worklog already exists
		exists := false
		for _, entry := range entries {
			if entry.JiraWorklogID != nil && *entry.JiraWorklogID == wl.ID {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		// Map Jira worklog to time entry
		entry := MapWorklogToTimeEntry(userID, wl)

		if err := s.timeEntryRepo.Create(ctx, entry); err != nil {
			return fmt.Errorf("creating time entry: %w", err)
		}
	}

	return nil
}

// SyncTimeEntryToJira logs a time entry to Jira
func (s *SyncService) SyncTimeEntryToJira(ctx context.Context, entry *timeentry.TimeEntry) error {
	if entry.JiraIssueKey == nil || *entry.JiraIssueKey == "" {
		return fmt.Errorf("no Jira issue key specified")
	}

	issueKey, timeSpentSeconds, comment, started := MapTimeEntryToJiraWorklog(entry)

	worklogID, err := s.client.LogWork(ctx, issueKey, timeSpentSeconds, started, comment)
	if err != nil {
		return fmt.Errorf("logging work to Jira: %w", err)
	}

	// Update the entry with Jira worklog ID
	entry.JiraWorklogID = &worklogID

	if err := s.timeEntryRepo.Update(ctx, entry); err != nil {
		return fmt.Errorf("updating time entry: %w", err)
	}

	return nil
}
