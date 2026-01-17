// Package invoice internal/domain/invoice/entity.go
package invoice

import (
	"github.com/google/uuid"

	"time"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusSent      Status = "sent"
	StatusPaid      Status = "paid"
	StatusOverdue   Status = "overdue"
	StatusCancelled Status = "cancelled"
)

type Invoice struct {
	ID            uuid.UUID `db:"id"`
	UserID        uuid.UUID `db:"user_id"`
	ClientID      uuid.UUID `db:"client_id"`
	InvoiceNumber string    `db:"invoice_number"`
	Status        Status    `db:"status"`
	IssueDate     time.Time `db:"issue_date"`
	DueDate       time.Time `db:"due_date"`
	Subtotal      float64   `db:"subtotal"`
	TaxRate       float64   `db:"tax_rate"`
	TaxAmount     float64   `db:"tax_amount"`
	Total         float64   `db:"total"`
	Currency      string    `db:"currency"`
	Notes         string    `db:"notes"`

	// Integration fields
	SquareInvoiceID *string `db:"square_invoice_id"`
	SquarePaymentID *string `db:"square_payment_id"`

	Items     []InvoiceItem
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type InvoiceItem struct {
	ID          uuid.UUID `db:"id"`
	InvoiceID   uuid.UUID `db:"invoice_id"`
	Description string    `db:"description"`
	Quantity    float64   `db:"quantity"`
	UnitPrice   float64   `db:"unit_price"`
	Amount      float64   `db:"amount"`
	SortOrder   int       `db:"sort_order"`
	CreatedAt   time.Time `db:"created_at"`
}

// Business logic methods
func (i *Invoice) CalculateTotals() {
	i.Subtotal = 0
	for _, item := range i.Items {
		i.Subtotal += item.Amount
	}
	i.TaxAmount = i.Subtotal * i.TaxRate / 100
	i.Total = i.Subtotal + i.TaxAmount
}

func (i *Invoice) MarkAsSent() error {
	if i.Status != StatusDraft {
		return ErrInvalidStatusTransition
	}
	i.Status = StatusSent
	return nil
}

func (i *Invoice) MarkAsPaid(paymentID string) error {
	if i.Status != StatusSent && i.Status != StatusOverdue {
		return ErrInvalidStatusTransition
	}
	i.Status = StatusPaid
	i.SquarePaymentID = &paymentID
	return nil
}
