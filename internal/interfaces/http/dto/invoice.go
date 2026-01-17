// internal/interfaces/http/dto/invoice.go
package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/invoice-app-be/internal/domain/invoice"
)

type CreateInvoiceRequest struct {
	ClientID  uuid.UUID              `json:"client_id" validate:"required"`
	IssueDate time.Time              `json:"issue_date" validate:"required"`
	DueDate   time.Time              `json:"due_date" validate:"required"`
	TaxRate   float64                `json:"tax_rate" validate:"gte=0,lte=100"`
	Currency  string                 `json:"currency" validate:"required,len=3"`
	Notes     string                 `json:"notes"`
	Items     []CreateInvoiceItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreateInvoiceItemDTO struct {
	Description string  `json:"description" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice   float64 `json:"unit_price" validate:"required,gte=0"`
}

type InvoiceResponse struct {
	ID            string           `json:"id"`
	ClientID      string           `json:"client_id"`
	InvoiceNumber string           `json:"invoice_number"`
	Status        string           `json:"status"`
	IssueDate     string           `json:"issue_date"`
	DueDate       string           `json:"due_date"`
	Subtotal      float64          `json:"subtotal"`
	TaxRate       float64          `json:"tax_rate"`
	TaxAmount     float64          `json:"tax_amount"`
	Total         float64          `json:"total"`
	Currency      string           `json:"currency"`
	Notes         string           `json:"notes"`
	Items         []InvoiceItemDTO `json:"items"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string           `json:"updated_at"`
}

type InvoiceItemDTO struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
}

func InvoiceFromDomain(inv *invoice.Invoice) InvoiceResponse {
	items := make([]InvoiceItemDTO, len(inv.Items))
	for i, item := range inv.Items {
		items[i] = InvoiceItemDTO{
			ID:          item.ID.String(),
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      item.Amount,
		}
	}

	return InvoiceResponse{
		ID:            inv.ID.String(),
		ClientID:      inv.ClientID.String(),
		InvoiceNumber: inv.InvoiceNumber,
		Status:        string(inv.Status),
		IssueDate:     inv.IssueDate.Format("2006-01-02"),
		DueDate:       inv.DueDate.Format("2006-01-02"),
		Subtotal:      inv.Subtotal,
		TaxRate:       inv.TaxRate,
		TaxAmount:     inv.TaxAmount,
		Total:         inv.Total,
		Currency:      inv.Currency,
		Notes:         inv.Notes,
		Items:         items,
		CreatedAt:     inv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     inv.UpdatedAt.Format(time.RFC3339),
	}
}
