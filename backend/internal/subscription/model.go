package subscription

import (
	"billingService/backend/internal/money"
	"billingService/backend/internal/plans"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	SubscriptionStatusTrial     Status = "trial"
	SubscriptionStatusActive    Status = "active"
	SubscriptionStatusPastDue   Status = "past_due"
	SubscriptionStatusCancelled Status = "cancelled"
)

// Base replicates gorm.Model
type Base struct {
	ID        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz" json:"updated_at"`
}

type Subscriptions struct {
	Base
	UserAccountID       uuid.UUID  `gorm:"column:user_account_id;type:uuid;not null" json:"user_account_id"`
	SubscriptionPlanID  int16      `gorm:"column:subscription_plan_id;type:smallint;not null" json:"subscription_plan_id"`
	TrialEndsAt         *time.Time `gorm:"column:trial_ends_at;type:timestamptz" json:"trial_ends_at"`
	CurrentPeriodEndsAt time.Time  `gorm:"column:current_period_ends_at;type:timestamptz;not null" json:"current_period_ends_at"`
	CancelAtPeriodEnd   bool       `gorm:"column:cancel_at_period_end;type:boolean;not null" json:"cancel_at_period_end"`
	Status              Status     `gorm:"column:status;type:subscription_status;not null" json:"status"`
	// Pointer allows the field to be nil until set
	CancelledAt *time.Time `gorm:"column:cancelled_at;type:timestamptz" json:"cancelled_at,omitempty"`
}

type SubscriptionsReadModel struct {
	Base
	UserAccountEmail    string                `gorm:"column:user_account_email;type:varchar;not null;unique" json:"user_account_email"`
	PlanName            string                `gorm:"column:plan_name;type:varchar;not null;unique" json:"plan_name"`
	PlanPrice           money.Money           `gorm:"embedded;embeddedPrefix:plan_price_" json:"plan_price"`
	PlanBillingInterval plans.BillingInterval `gorm:"column:plan_billing_interval;type:billing_interval;not null" json:"subscription_plan_billing_interval"`
	TrialEndsAt         *time.Time            `gorm:"column:trial_ends_at;type:timestamptz" json:"trial_ends_at"`
	CurrentPeriodEndsAt time.Time             `gorm:"column:current_period_ends_at;type:timestamptz;not null" json:"current_period_ends_at"`
	CancelAtPeriodEnd   bool                  `gorm:"column:cancel_at_period_end;type:boolean;not null" json:"cancel_at_period_end"`
	Status              Status                `gorm:"column:status;type:subscription_status;not null" json:"status"`
	CancelledAt         *time.Time            `gorm:"column:cancelled_at;type:timestamptz" json:"cancelled_at,omitempty"`
}

func (Subscriptions) TableName() string {
	return "subscriptions"
}

func (subscription *Subscriptions) Activate() error {
	if subscription.Status != SubscriptionStatusTrial && subscription.Status != SubscriptionStatusPastDue {
		return fmt.Errorf("cannot activate a subscription with status: %s", subscription.Status)
	}
	subscription.Status = SubscriptionStatusActive
	return nil
}

func (subscription *Subscriptions) Cancel() error {
	if subscription.Status == SubscriptionStatusCancelled {
		return fmt.Errorf("cannot cancel a plan with status: %s", subscription.Status)
	}
	subscription.Status = SubscriptionStatusCancelled
	return nil
}

func (subscription *Subscriptions) EnterGracePeriod() error {
	if subscription.Status != SubscriptionStatusActive {
		return fmt.Errorf("cannot enter grace period with status: %s", subscription.Status)
	}
	subscription.Status = SubscriptionStatusPastDue
	return nil
}
