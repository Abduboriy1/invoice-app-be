// internal/infrastructure/pdf/generator.go
package pdf

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jung-kurt/gofpdf"

	"github.com/invoice-app-be/internal/domain/invoice"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(ctx context.Context, inv *invoice.Invoice) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 24)
	pdf.Cell(40, 10, "INVOICE")
	pdf.Ln(15)

	// Invoice details
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice #: %s", inv.InvoiceNumber))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", inv.IssueDate.Format("2006-01-02")))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Due: %s", inv.DueDate.Format("2006-01-02")))
	pdf.Ln(15)

	// Items
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(100, 10, "Description")
	pdf.Cell(30, 10, "Quantity")
	pdf.Cell(40, 10, "Amount")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	for _, item := range inv.Items {
		pdf.Cell(100, 8, item.Description)
		pdf.Cell(30, 8, fmt.Sprintf("%.2f", item.Quantity))
		pdf.Cell(40, 8, fmt.Sprintf("$%.2f", item.Amount))
		pdf.Ln(8)
	}

	// Total
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(130, 10, "Total:")
	pdf.Cell(40, 10, fmt.Sprintf("$%.2f", inv.Total))

	// Convert to bytes - FIXED
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("generating PDF: %w", err)
	}

	return buf.Bytes(), nil
}
