// internal/domain/timeentry/service.go
package timeentry

import (
	"context"
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
		return nil, err
	}

	return entry, nil
}

type CreateTimeEntryRequest struct {
	Description string
	Hours       float64
	Date        time.Time
	IsBillable  bool
}
