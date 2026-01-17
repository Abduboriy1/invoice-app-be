// internal/interfaces/http/handlers/invoice.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/invoice-app-be/internal/domain/invoice"
	"github.com/invoice-app-be/internal/interfaces/http/dto"
	"github.com/invoice-app-be/internal/interfaces/http/middleware"
)

var validate = validator.New()

type InvoiceHandler struct {
	service    *invoice.Service
	clientRepo interface{} // Placeholder for now
}

func NewInvoiceHandler(service *invoice.Service, clientRepo interface{}) *InvoiceHandler {
	return &InvoiceHandler{
		service:    service,
		clientRepo: clientRepo,
	}
}

func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	var req dto.CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Map DTO to domain request
	domainReq := invoice.CreateInvoiceRequest{
		ClientID:  req.ClientID,
		IssueDate: req.IssueDate,
		DueDate:   req.DueDate,
		TaxRate:   req.TaxRate,
		Currency:  req.Currency,
		Notes:     req.Notes,
		Items:     make([]invoice.CreateInvoiceItemRequest, len(req.Items)),
	}

	for i, item := range req.Items {
		domainReq.Items[i] = invoice.CreateInvoiceItemRequest{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}

	inv, err := h.service.CreateInvoice(r.Context(), userID, domainReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	respondJSON(w, http.StatusCreated, dto.InvoiceFromDomain(inv))
}

func (h *InvoiceHandler) List(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, []dto.InvoiceResponse{})
}

func (h *InvoiceHandler) Get(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *InvoiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *InvoiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusNoContent, nil)
}

func (h *InvoiceHandler) Send(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

func (h *InvoiceHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	invoiceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	pdfBytes, err := h.service.GeneratePDF(r.Context(), userID, invoiceID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate PDF")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=invoice.pdf")
	w.Write(pdfBytes)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
