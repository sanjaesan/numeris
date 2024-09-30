package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/numeris/internal/domain"
	"gorm.io/gorm"
)

// InvoiceRepository -
type InvoiceRepository interface {
	InvoiceDB
}

// InvoiceDB -
type InvoiceDB interface {
	ByID(ctx context.Context, id uint) (*domain.Invoice, error)
	ByInvoiceNo(ctx context.Context, invoiceNo string) (*domain.Invoice, error)

	Create(ctx context.Context, invoice *domain.Invoice) error
	Update(ctx context.Context, s *domain.Invoice) error
	Delete(ctx context.Context, userID uint) error
}

type invoiceRepository struct {
	InvoiceDB
}

// NewInvoiceRepository -
func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	s := newInvoiceValidator(&invoiceGorm{db}, uuid.New())
	return &invoiceRepository{
		InvoiceDB: s,
	}
}

type invoiceValidatorFunc func(*domain.Invoice) error

func runInvoiceValidatorFuncs(project *domain.Invoice, fns ...invoiceValidatorFunc) error {
	for _, fn := range fns {
		if err := fn(project); err != nil {
			return err
		}
	}
	return nil
}

var _ InvoiceDB = &invoiceValidator{}

type invoiceValidator struct {
	InvoiceDB
	uuid uuid.UUID
}

func newInvoiceValidator(idb InvoiceDB, uuid uuid.UUID) *invoiceValidator {
	return &invoiceValidator{
		InvoiceDB: idb,
		uuid:      uuid,
	}
}

func (s *invoiceValidator) Create(ctx context.Context, invoice *domain.Invoice) error {
	err := runInvoiceValidatorFuncs(invoice,
		s.invoiceNumber, s.dateRequired, s.billingCurrency,
		s.emailRequired, s.nameRequired, s.itemsRequired)
	if err != nil {
		return err
	}
	return s.InvoiceDB.Create(ctx, invoice)
}

func (s *invoiceValidator) Update(ctx context.Context, invoice *domain.Invoice) error {
	err := runInvoiceValidatorFuncs(invoice,
		s.dateRequired, s.billingCurrency,
		s.emailRequired, s.nameRequired, s.itemsRequired)
	if err != nil {
		return err
	}
	return s.InvoiceDB.Update(ctx, invoice)
}

func (s *invoiceValidator) Delete(ctx context.Context, id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return s.InvoiceDB.Delete(ctx, id)
}
func (s *invoiceValidator) invoiceNumber(inv *domain.Invoice) error {
	if inv.InvoiceNo == "" {
		return ErrInvoiceNoRequired
	}
	inv.InvoiceNo = fmt.Sprintf("Invoice-%s", s.uuid.String()[:13])
	return nil
}

func (s *invoiceValidator) dateRequired(inv *domain.Invoice) error {
	if inv.IssueDate.String() == "" || inv.DueDate.String() == "" {
		return ErrDateRequired
	}
	return nil
}
func (s *invoiceValidator) billingCurrency(inv *domain.Invoice) error {
	if inv.BillingCurrency == "" {
		return ErrCurrencyRequired
	}
	return nil
}
func (s *invoiceValidator) emailRequired(inv *domain.Invoice) error {
	if inv.SenderDetails.Email == "" || inv.CustomerDetails.Email == "" {
		return ErrEmailRequired
	}
	return nil
}
func (s *invoiceValidator) nameRequired(inv *domain.Invoice) error {
	if inv.SenderDetails.Name == "" || inv.CustomerDetails.Name == "" {
		return ErrNameRequired
	}
	return nil
}
func (s *invoiceValidator) itemsRequired(inv *domain.Invoice) error {
	if len(inv.InvoiceItems) < 1 {
		return ErrItemsRequired
	}
	return nil
}

type invoiceGorm struct {
	db *gorm.DB
}

func (s *invoiceGorm) ByID(ctx context.Context, id uint) (*domain.Invoice, error) {
	var invoice domain.Invoice
	db := s.db.Where("id = ?", id)
	err := first(db, ctx, &invoice)
	return &invoice, err
}
func (s *invoiceGorm) ByInvoiceNo(ctx context.Context, invoice_no string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	db := s.db.Where("invoice_no = ?", invoice_no)
	err := first(db, ctx, &invoice)
	return &invoice, err
}

// Create -
func (s *invoiceGorm) Create(ctx context.Context, invoice *domain.Invoice) error {
	return s.db.Create(invoice).Error
}

// Update -
func (s *invoiceGorm) Update(ctx context.Context, invoice *domain.Invoice) error {
	return s.db.WithContext(ctx).Save(invoice).Error
}

// Delete will permanently delete the invoice with the provided ID
func (s *invoiceGorm) Delete(ctx context.Context, id uint) error {
	inv := domain.Invoice{Model: gorm.Model{ID: id}}
	return s.db.WithContext(ctx).Unscoped().Delete(&inv).Error
}

// first will query using the provided gorm.DB and it will
// get the first item returned and place it into dst. If
// nothing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, ctx context.Context, dst interface{}) error {
	err := db.WithContext(ctx).First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
