package domain

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	InvoiceNo          string `gorm:"not null;unique_index"`
	IssueDate          time.Time
	DueDate            time.Time
	BillingCurrency    string
	SenderDetails      SenderDetails
	CustomerDetails    CustomerDetails
	InvoiceItems       []InvoiceItem
	SubTotal           float64
	Discount           float64
	Total              float64
	PaymentInformation *PaymentInformation
}

type InvoiceItem struct {
	gorm.Model
	InvoiceID    uint64 `gorm:"not null"`
	Description  string `gorm:"not null"`
	Unit         string
	PricePerUnit float64
	Total        float64
}

type SenderDetails struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Address string
	Email   string `gorm:"not null"`
	Phone   string
}

type CustomerDetails struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
	Phone string
}

type PaymentInformation struct {
	gorm.Model
	InvoiceID   uint64 `gorm:"not null"`
	AccountName string
	AccountNo   string
	RoutingNo   string
	BankName    string
	Address     string
}
