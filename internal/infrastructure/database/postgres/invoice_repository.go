// internal/infrastructure/database/postgres/invoice_repository.go
package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/invoice-app-be/internal/domain/invoice"
)

type InvoiceRepository struct {
	db *sqlx.DB
}

func NewInvoiceRepository(db *sqlx.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(ctx context.Context, inv *invoice.Invoice) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert invoice
	query := `
        INSERT INTO invoices (id, user_id, client_id, invoice_number, status, issue_date, due_date, 
                            subtotal, tax_rate, tax_amount, total, currency, notes, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `
	_, err = tx.ExecContext(ctx, query, inv.ID, inv.UserID, inv.ClientID, inv.InvoiceNumber, inv.Status,
		inv.IssueDate, inv.DueDate, inv.Subtotal, inv.TaxRate, inv.TaxAmount, inv.Total, inv.Currency,
		inv.Notes, inv.CreatedAt, inv.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert items
	for _, item := range inv.Items {
		itemQuery := `
            INSERT INTO invoice_items (id, invoice_id, description, quantity, unit_price, amount, sort_order, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        `
		_, err = tx.ExecContext(ctx, itemQuery, item.ID, item.InvoiceID, item.Description, item.Quantity,
			item.UnitPrice, item.Amount, item.SortOrder, item.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *InvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*invoice.Invoice, error) {
	var inv invoice.Invoice
	query := `
        SELECT id, user_id, client_id, invoice_number, status, issue_date, due_date,
               subtotal, tax_rate, tax_amount, total, currency, notes, created_at, updated_at
        FROM invoices WHERE id = $1
    `
	if err := r.db.GetContext(ctx, &inv, query, id); err != nil {
		return nil, fmt.Errorf("getting invoice: %w", err)
	}

	// Get items
	var items []invoice.InvoiceItem
	itemQuery := `SELECT id, invoice_id, description, quantity, unit_price, amount, sort_order, created_at 
                  FROM invoice_items WHERE invoice_id = $1 ORDER BY sort_order`
	if err := r.db.SelectContext(ctx, &items, itemQuery, id); err != nil {
		return nil, fmt.Errorf("getting invoice items: %w", err)
	}
	inv.Items = items

	return &inv, nil
}

func (r *InvoiceRepository) GetByUserID(ctx context.Context, userID uuid.UUID, filters invoice.ListFilters) ([]invoice.Invoice, error) {
	query := `
        SELECT id, user_id, client_id, invoice_number, status, issue_date, due_date,
               subtotal, tax_rate, tax_amount, total, currency, notes, created_at, updated_at
        FROM invoices WHERE user_id = $1 ORDER BY created_at DESC
    `
	var invoices []invoice.Invoice
	if err := r.db.SelectContext(ctx, &invoices, query, userID); err != nil {
		return nil, fmt.Errorf("getting invoices: %w", err)
	}
	return invoices, nil
}

func (r *InvoiceRepository) Update(ctx context.Context, inv *invoice.Invoice) error {
	query := `
        UPDATE invoices SET status = $2, subtotal = $3, tax_rate = $4, tax_amount = $5, 
                          total = $6, notes = $7, updated_at = $8
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, inv.ID, inv.Status, inv.Subtotal, inv.TaxRate,
		inv.TaxAmount, inv.Total, inv.Notes, inv.UpdatedAt)
	return err
}

func (r *InvoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM invoices WHERE id = $1", id)
	return err
}

func (r *InvoiceRepository) GetNextInvoiceNumber(ctx context.Context, userID uuid.UUID) (string, error) {
	var count int
	query := `SELECT COUNT(*) FROM invoices WHERE user_id = $1`
	if err := r.db.GetContext(ctx, &count, query, userID); err != nil {
		return "", err
	}
	return fmt.Sprintf("INV-%05d", count+1), nil
}
