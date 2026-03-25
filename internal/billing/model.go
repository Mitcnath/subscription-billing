package billing

import (
	"billingService/backend/internal/accounts"
	"billingService/backend/internal/statuses"
	"time"

	"github.com/google/uuid"
)

type BillingStatus string

const (
	// The invoice has been paid. The invoice will be closed and no further payment attempts will be made.
	BillingStatusPaid BillingStatus = "paid"
	// The invoice is still open and payment has not yet been attempted.
	BillingStatusOpen BillingStatus = "open"
	// The invoice was cancelled before it was ever paid.
	BillingStatusVoid BillingStatus = "void"
	// Payment was attempted and failed, and we have stopped retrying the payment.
	BillingStatusUncollectible BillingStatus = "uncollectible"
)

func (billingStatus BillingStatus) Valid() bool {
	switch billingStatus {
	case BillingStatusPaid, BillingStatusOpen, BillingStatusVoid, BillingStatusUncollectible:
		return true
	}
	return false
}

type Base struct {
	ID        int64     `gorm:"column:id;type:bigint;not null;primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null" json:"updated_at"`
}

type BillingHistories struct {
	Base
	UserAccountID  uuid.UUID              `gorm:"column:user_account_id;type:uuid;not null" json:"user_account_id"`
	UserAccount    accounts.UserAccounts  `gorm:"foreignKey:user_account_id"`
	SubscriptionID int64                  `gorm:"column:subscription_id;type:bigint;not null" json:"subscription_id"`
	Subscription   statuses.Subscriptions `gorm:"foreignKey:subscription_id"`
	Status         BillingStatus          `gorm:"column:status;type:billing_status;not null" json:"status"`
	AmountPaid     int64                  `gorm:"column:amount_paid;type:bigint;not null" json:"amount_paid"`
	PdfURL         string                 `gorm:"column:pdf_url;type:text;not null" json:"pdf_url"`
}

func (BillingHistories) TableName() string {
	return "billing_histories"
}
