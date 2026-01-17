// internal/domain/invoice/types.go
package invoice

import (
	"time"

	"github.com/google/uuid"
)

type CreateInvoiceRequest struct {
	ClientID  uuid.UUID
	IssueDate time.Time
	DueDate   time.Time
	TaxRate   float64
	Currency  string
	Notes     string
	Items     []CreateInvoiceItemRequest
}

type CreateInvoiceItemRequest struct {
	Description string
	Quantity    float64
	UnitPrice   float64
}
