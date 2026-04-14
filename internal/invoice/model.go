package invoice

import (
	"billingService/backend/internal/money"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	// The invoice has been paid. The invoice will be closed and no further payment attempts will be made.
	InvoiceStatusPaid Status = "paid"
	// The invoice is still open and payment has not yet been attempted.
	InvoiceStatusOpen Status = "open"
	// The invoice was cancelled before it was ever paid.
	InvoiceStatusVoid Status = "void"
	// Payment was attempted and failed, and we have stopped retrying the payment.
	InvoiceStatusUncollectible Status = "uncollectible"
)

func (status Status) Valid() bool {
	switch status {
	case InvoiceStatusPaid, InvoiceStatusOpen, InvoiceStatusVoid, InvoiceStatusUncollectible:
		return true
	}
	return false
}

type Base struct {
	ID        int64     `gorm:"column:id;type:bigint;not null;primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null" json:"updated_at"`
}

type Invoice struct {
	Base
	UserAccountID  uuid.UUID   `gorm:"column:user_account_id;type:uuid;not null" json:"user_account_id"`
	SubscriptionID int64       `gorm:"column:subscription_id;type:bigint;not null" json:"subscription_id"`
	Status         Status      `gorm:"column:status;type:invoice_status;not null" json:"status"`
	Paid           money.Money `gorm:"embedded"`
	PdfURL         string      `gorm:"column:pdf_url;type:text;not null" json:"pdf_url"`
}

func (Invoice) TableName() string {
	return "invoices"
}

type InvoiceReadModel struct {
	Base
	UserEmail         string      `gorm:"column:user_email;type:varchar;not null" json:"user_email"`
	SubscriptionName  string      `gorm:"column:subscription_name;type:varchar;not null" json:"subscription_name"`
	SubscriptionPrice money.Money `gorm:"embedded;embeddedPrefix:subscription_price_" json:"subscription_price"`
	Status            Status      `gorm:"column:status;type:invoice_status;not null" json:"status"`
	Paid              money.Money `gorm:"embedded;embeddedPrefix:paid_" json:"paid"`
	PdfURL            string      `gorm:"column:pdf_url;type:text;not null" json:"pdf_url"`
}

func (invoice *Invoice) MarkPaid(money money.Money) error {
	if invoice.Status != InvoiceStatusOpen {
		return fmt.Errorf("cannot mark invoice as paid with status: %s", invoice.Status)
	}
	invoice.Status = InvoiceStatusPaid
	invoice.Paid.Amount = money.Amount
	return nil
}

func (invoice *Invoice) Void() error {
	if invoice.Status != InvoiceStatusOpen {
		return fmt.Errorf("cannot void an invoice with status: %s", invoice.Status)
	}
	invoice.Status = InvoiceStatusVoid
	return nil
}
