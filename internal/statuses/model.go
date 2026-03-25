package statuses

import (
	"billingService/backend/internal/accounts"
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusTrial    SubscriptionStatus = "trial"
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusPastDue  SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled SubscriptionStatus = "cancelled"
)

// Base replicates gorm.Model
type Base struct {
	ID        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz" json:"updated_at"`
}

type Subscriptions struct {
	Base
	UserAccountID       uuid.UUID             `gorm:"column:user_account_id;type:uuid;not null" json:"user_account_id"`
	UserAccount         accounts.UserAccounts `gorm:"foreignKey:user_account_id"`
	SubscriptionPlanID  int16                 `gorm:"column:subscription_plan_id;type:smallint;not null" json:"subscription_plan_id"`
	TrialEndsAt         *time.Time            `gorm:"column:trial_ends_at;type:timestamptz" json:"trial_ends_at"`
	CurrentPeriodEndsAt time.Time             `gorm:"column:current_period_ends_at;type:timestamptz;not null" json:"current_period_ends_at"`
	CancelAtPeriodEnd   bool                  `gorm:"column:cancel_at_period_end;type:boolean;not null" json:"cancel_at_period_end"`
	Status              SubscriptionStatus    `gorm:"column:status;type:subscription_status;not null" json:"status"`
	// Pointer allows the field to be nil until set
	CancelledAt *time.Time `gorm:"column:cancelled_at;type:timestamptz" json:"cancelled_at,omitempty"`
}

func (Subscriptions) TableName() string {
	return "subscriptions"
}
