// internal/infrastructure/integrations/jira/sync.go
package jira

import (
	"context"
	"fmt"
	"time"

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

// PullWorklogsByDateRange pulls worklogs from Jira for a date range
func (s *SyncService) PullWorklogsByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, issueKeys []string) (int, error) {
	if s.client == nil {
		return 0, fmt.Errorf("Jira client not configured")
	}

	var worklogs []Worklog
	var err error

	if len(issueKeys) > 0 {
		// Use provided issue keys
		worklogs, err = s.client.GetWorklogsByDateRange(ctx, issueKeys, startDate, endDate)
		if err != nil {
			return 0, fmt.Errorf("fetching worklogs: %w", err)
		}
	} else {
		// Get all issues with worklogs in the date range using new JQL API
		fmt.Printf("Searching for issues with worklogs between %s and %s\n",
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"))

		keys, err := s.client.GetAllIssuesWithWorklogsByDateRange(ctx, startDate, endDate)
		if err != nil {
			return 0, fmt.Errorf("searching for issues with worklogs: %w", err)
		}

		fmt.Printf("Found %d issues with worklogs\n", len(keys))

		if len(keys) == 0 {
			return 0, nil // No issues found
		}

		worklogs, err = s.client.GetWorklogsByDateRange(ctx, keys, startDate, endDate)
		if err != nil {
			return 0, fmt.Errorf("fetching worklogs: %w", err)
		}
	}

	fmt.Printf("Found %d worklogs to process\n", len(worklogs))
	fmt.Println(worklogs)

	// Create time entries from worklogs
	count := 0
	for _, wl := range worklogs {
		// Check if already synced
		entries, _ := s.timeEntryRepo.GetByUserID(ctx, userID, "", "")
		exists := false
		for _, entry := range entries {
			if entry.JiraWorklogID != nil && *entry.JiraWorklogID == wl.ID {
				exists = true
				break
			}
		}

		if exists {
			fmt.Printf("Skipping worklog %s (already synced)\n", wl.ID)
			continue
		}

		// Create new time entry
		entry := MapWorklogToTimeEntry(userID, wl)
		if err := s.timeEntryRepo.Create(ctx, entry); err != nil {
			return count, fmt.Errorf("creating time entry: %w", err)
		}
		fmt.Printf("Created time entry for worklog %s\n", wl.ID)
		count++
	}

	return count, nil
}

// SyncWorklogsForIssue syncs all worklogs for a specific issue
func (s *SyncService) SyncWorklogsForIssue(ctx context.Context, userID uuid.UUID, issueKey string) error {
	if s.client == nil {
		return fmt.Errorf("Jira client not configured")
	}

	worklogs, err := s.client.GetWorklogs(ctx, issueKey)
	if err != nil {
		return fmt.Errorf("fetching worklogs: %w", err)
	}

	for _, wl := range worklogs {
		// Check if already synced
		entries, _ := s.timeEntryRepo.GetByUserID(ctx, userID, "", "")
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

		// Create time entry
		entry := MapWorklogToTimeEntry(userID, wl)
		if err := s.timeEntryRepo.Create(ctx, entry); err != nil {
			return fmt.Errorf("creating time entry: %w", err)
		}
	}

	return nil
}

// PushTimeEntryToJira pushes a time entry to Jira
func (s *SyncService) PushTimeEntryToJira(ctx context.Context, entry *timeentry.TimeEntry, issueKey string) error {
	if s.client == nil {
		return fmt.Errorf("Jira client not configured")
	}

	// Convert hours to seconds
	timeSpentSeconds := int(entry.Hours * 3600)

	// Log work to Jira
	worklogID, err := s.client.LogWork(ctx, issueKey, timeSpentSeconds, entry.Date, entry.Description)
	if err != nil {
		return fmt.Errorf("logging work to Jira: %w", err)
	}

	// Update entry with Jira info
	entry.JiraIssueKey = &issueKey
	entry.JiraWorklogID = &worklogID
	now := time.Now()
	entry.JiraSyncedAt = &now
	entry.UpdatedAt = now

	if err := s.timeEntryRepo.Update(ctx, entry); err != nil {
		return fmt.Errorf("updating time entry: %w", err)
	}

	return nil
}
