package square

import (
	"context"

	"github.com/invoice-app-be/internal/domain/invoice"
)

type InvoiceService struct {
	client *Client
}

func NewInvoiceService(client *Client) *InvoiceService {
	return &InvoiceService{client: client}
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, inv *invoice.Invoice) (string, error) {
	// TODO: Implement Square invoice creation
	return "", nil
}

func (s *InvoiceService) SendInvoice(ctx context.Context, squareInvoiceID string) error {
	// TODO: Implement
	return nil
}
