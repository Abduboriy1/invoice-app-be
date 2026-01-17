// internal/domain/timeentry/service.go
package timeentry

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JiraClient interface {
	LogWork(ctx context.Context, issueKey string, timeSpentSeconds int, started time.Time, comment string) (string, error)
}

type Service struct {
	repo       Repository
	jiraClient JiraClient
}

func NewService(repo Repository, jiraClient JiraClient) *Service {
	return &Service{
		repo:       repo,
		jiraClient: jiraClient,
	}
}

type CreateTimeEntryRequest struct {
	Description string
	Hours       float64
	Date        time.Time
	IsBillable  bool
}

type UpdateTimeEntryRequest struct {
	Description string
	Hours       float64
	Date        time.Time
	IsBillable  bool
}

func (s *Service) CreateTimeEntry(ctx context.Context, userID uuid.UUID, req CreateTimeEntryRequest) (*TimeEntry, error) {
	entry := &TimeEntry{
		ID:          uuid.New(),
		UserID:      userID,
		Description: req.Description,
		Hours:       req.Hours,
		Date:        req.Date,
		IsBillable:  req.IsBillable,
		IsInvoiced:  false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("creating time entry: %w", err)
	}

	return entry, nil
}

func (s *Service) ListTimeEntries(ctx context.Context, userID uuid.UUID) ([]TimeEntry, error) {
	entries, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listing time entries: %w", err)
	}
	return entries, nil
}

func (s *Service) GetTimeEntry(ctx context.Context, userID, entryID uuid.UUID) (*TimeEntry, error) {
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return nil, fmt.Errorf("getting time entry: %w", err)
	}

	if entry.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	return entry, nil
}

func (s *Service) UpdateTimeEntry(ctx context.Context, userID, entryID uuid.UUID, req UpdateTimeEntryRequest) (*TimeEntry, error) {
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return nil, fmt.Errorf("getting time entry: %w", err)
	}

	if entry.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	entry.Description = req.Description
	entry.Hours = req.Hours
	entry.Date = req.Date
	entry.IsBillable = req.IsBillable
	entry.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, fmt.Errorf("updating time entry: %w", err)
	}

	return entry, nil
}

func (s *Service) DeleteTimeEntry(ctx context.Context, userID, entryID uuid.UUID) error {
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return fmt.Errorf("getting time entry: %w", err)
	}

	if entry.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	if err := s.repo.Delete(ctx, entryID); err != nil {
		return fmt.Errorf("deleting time entry: %w", err)
	}

	return nil
}

func (s *Service) SyncToJira(ctx context.Context, userID, entryID uuid.UUID, issueKey string) error {
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return fmt.Errorf("getting time entry: %w", err)
	}

	if entry.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	if s.jiraClient == nil {
		return fmt.Errorf("Jira integration not configured")
	}

	// Convert hours to seconds
	timeSpentSeconds := int(entry.Hours * 3600)

	// Log work to Jira
	worklogID, err := s.jiraClient.LogWork(ctx, issueKey, timeSpentSeconds, entry.Date, entry.Description)
	if err != nil {
		return fmt.Errorf("logging work to Jira: %w", err)
	}

	// Update entry with Jira info
	entry.JiraIssueKey = &issueKey
	entry.JiraWorklogID = &worklogID
	now := time.Now()
	entry.JiraSyncedAt = &now
	entry.UpdatedAt = now

	if err := s.repo.Update(ctx, entry); err != nil {
		return fmt.Errorf("updating time entry: %w", err)
	}

	return nil
}
