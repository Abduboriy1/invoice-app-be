// internal/domain/invoice/service.go
package invoice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvoiceNotFound         = fmt.Errorf("invoice not found")
	ErrInvalidStatusTransition = fmt.Errorf("invalid status transition")
	ErrUnauthorized            = fmt.Errorf("unauthorized access")
)

type Service struct {
	repo      Repository
	pdfGen    PDFGenerator
	squareAPI SquareAPI
}

func NewService(repo Repository, pdfGen PDFGenerator, squareAPI SquareAPI) *Service {
	return &Service{
		repo:      repo,
		pdfGen:    pdfGen,
		squareAPI: squareAPI,
	}
}

func (s *Service) CreateInvoice(ctx context.Context, userID uuid.UUID, req CreateInvoiceRequest) (*Invoice, error) {
	// Generate invoice number
	invoiceNum, err := s.repo.GetNextInvoiceNumber(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("generating invoice number: %w", err)
	}

	invoice := &Invoice{
		ID:            uuid.New(),
		UserID:        userID,
		ClientID:      req.ClientID,
		InvoiceNumber: invoiceNum,
		Status:        StatusDraft,
		IssueDate:     req.IssueDate,
		DueDate:       req.DueDate,
		TaxRate:       req.TaxRate,
		Currency:      req.Currency,
		Notes:         req.Notes,
		Items:         make([]InvoiceItem, len(req.Items)),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Add items
	for i, item := range req.Items {
		invoice.Items[i] = InvoiceItem{
			ID:          uuid.New(),
			InvoiceID:   invoice.ID,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      item.Quantity * item.UnitPrice,
			SortOrder:   i,
			CreatedAt:   time.Now(),
		}
	}

	invoice.CalculateTotals()

	if err := s.repo.Create(ctx, invoice); err != nil {
		return nil, fmt.Errorf("creating invoice: %w", err)
	}

	return invoice, nil
}

func (s *Service) SendInvoice(ctx context.Context, userID, invoiceID uuid.UUID) error {
	invoice, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	if invoice.UserID != userID {
		return ErrUnauthorized
	}

	if err := invoice.MarkAsSent(); err != nil {
		return err
	}

	// Optionally sync to Square
	if s.squareAPI != nil {
		squareID, err := s.squareAPI.CreateInvoice(ctx, invoice)
		if err != nil {
			// Log but don't fail - Square is optional
			// logger.Error("failed to sync to Square", "error", err)
		} else {
			invoice.SquareInvoiceID = &squareID
		}
	}

	return s.repo.Update(ctx, invoice)
}

func (s *Service) GeneratePDF(ctx context.Context, userID, invoiceID uuid.UUID) ([]byte, error) {
	invoice, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if invoice.UserID != userID {
		return nil, ErrUnauthorized
	}

	return s.pdfGen.Generate(ctx, invoice)
}

// Interfaces for dependencies (ports)
type PDFGenerator interface {
	Generate(ctx context.Context, invoice *Invoice) ([]byte, error)
}

type SquareAPI interface {
	CreateInvoice(ctx context.Context, invoice *Invoice) (string, error)
	GetPaymentStatus(ctx context.Context, invoiceID string) (string, error)
}
