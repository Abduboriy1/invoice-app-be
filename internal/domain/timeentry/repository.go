// internal/domain/timeentry/repository.go
package timeentry

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, entry *TimeEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*TimeEntry, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]TimeEntry, error)
	GetByJiraWorklogID(ctx context.Context, worklogID string) (*TimeEntry, error)
	Update(ctx context.Context, entry *TimeEntry) error
	Delete(ctx context.Context, id uuid.UUID) error
}
