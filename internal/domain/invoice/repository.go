// internal/domain/invoice/repository.go
package invoice

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the contract for invoice persistence
type Repository interface {
	Create(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*Invoice, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, filters ListFilters) ([]Invoice, error)
	Update(ctx context.Context, invoice *Invoice) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetNextInvoiceNumber(ctx context.Context, userID uuid.UUID) (string, error)
}

type ListFilters struct {
	Status   *Status
	ClientID *uuid.UUID
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}
