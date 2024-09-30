package service

import (
	"github.com/numeris/internal/domain"
	"github.com/numeris/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InvoiceService -
type InvoiceService struct {
	repository.InvoiceRepository
	db *gorm.DB
}

// ServiceConfig -
type ServiceConfig func(*InvoiceService) error

// NewInvoiceService -
func NewInvoiceService(cfgs ...ServiceConfig) (*InvoiceService, error) {
	var s InvoiceService
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// WithGorm -
func WithGorm(dialect gorm.Dialector, config *gorm.Config) ServiceConfig {
	return func(s *InvoiceService) error {
		db, err := gorm.Open(dialect, config)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

// WithLogMode -
func WithLogMode() ServiceConfig {
	return func(s *InvoiceService) error {
		s.db.Logger.LogMode(logger.Silent)
		return nil
	}
}

// WithInvoice -
func WithInvoice() ServiceConfig {
	return func(s *InvoiceService) error {
		s.InvoiceRepository = repository.NewInvoiceRepository(s.db)
		// s.Invoice = repository.NewInvoiceRepository(s.db)
		return nil
	}
}

// Close -
func (s *InvoiceService) Close() error {
	sqlDB, err := s.db.DB()
	if err == nil {
		return sqlDB.Close()
	}
	return err
}

// DestructiveReset drops the all tables and rebuilds them
func (s *InvoiceService) DestructiveReset() error {
	err := s.db.Migrator().DropTable(&domain.Invoice{}, &domain.InvoiceItem{},
		&domain.PaymentInformation{})
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate all tables
func (s *InvoiceService) AutoMigrate() error {
	return s.db.AutoMigrate(&domain.Invoice{}, &domain.InvoiceItem{},
		&domain.PaymentInformation{})
}
