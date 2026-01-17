// internal/infrastructure/integrations/square/client.go
package square

import (
	"context"

	"github.com/invoice-app-be/internal/domain/invoice"
)

type Client struct {
	accessToken string
	environment string
}

func NewClient(accessToken, environment string) *Client {
	return &Client{
		accessToken: accessToken,
		environment: environment,
	}
}

func (c *Client) CreateInvoice(ctx context.Context, inv *invoice.Invoice) (string, error) {
	// TODO: Implement Square API integration
	return "square-invoice-id", nil
}

func (c *Client) GetPaymentStatus(ctx context.Context, invoiceID string) (string, error) {
	// TODO: Implement Square API integration
	return "paid", nil
}
