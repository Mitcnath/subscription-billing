package plans

import "time"

type BillingInterval string

const (
	BillingIntervalDaily      BillingInterval = "daily"
	BillingIntervalWeekly     BillingInterval = "weekly"
	BillingIntervalBiWeekly   BillingInterval = "bi_weekly"
	BillingIntervalMonthly    BillingInterval = "monthly"
	BillingIntervalQuarterly  BillingInterval = "quarterly"
	BillingIntervalSemiAnnual BillingInterval = "semi_annual"
	BillingIntervalAnnual     BillingInterval = "annual"
)

// Validation function for BillingInterval type
func (billingInterval BillingInterval) Valid() bool {
	switch billingInterval {
	case BillingIntervalDaily, BillingIntervalWeekly, BillingIntervalBiWeekly,
		BillingIntervalMonthly, BillingIntervalQuarterly, BillingIntervalSemiAnnual, BillingIntervalAnnual:
		return true
	}
	return false
}

type PlanStatus string

const (
	PlanStatusActive     PlanStatus = "active"
	PlanStatusDeprecated PlanStatus = "deprecated"
)

func (planStatus PlanStatus) Valid() bool {
	switch planStatus {
	case PlanStatusActive, PlanStatusDeprecated:
		return true
	}
	return false
}

// Base replicates gorm.Model
type Base struct {
	ID        int16     `gorm:"column:id;type:smallint;not null;primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null" json:"updated_at"`
}

type SubscriptionPlans struct {
	Base
	Name            string          `gorm:"column:name;type:varchar;not null;unique" json:"name"`
	Amount          uint64          `gorm:"column:amount;type:bigint;not null;check:amount>0" json:"amount"`
	Currency        string          `gorm:"column:currency;type:varchar;not null" json:"currency"`
	Description     string          `gorm:"column:description;type:text;not null;default:''" json:"description"`
	BillingInterval BillingInterval `gorm:"column:billing_interval;type:billing_interval;not null" json:"billing_interval"`
	Status          PlanStatus      `gorm:"column:status;type:plan_status;not null;default:'active'" json:"status"`
}

func (SubscriptionPlans) TableName() string {
	return "subscription_plans"
}
