package payment

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        int64     `gorm:"column:id;type:bigint;not null;primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null" json:"updated_at"`
}

type Method struct {
	Base
	UserAccountID uuid.UUID `gorm:"column:user_account_id;type:uuid;not null" json:"user_account_id"`
	ExternalID    string    `gorm:"column:external_id;type:varchar;not null" json:"external_id"`
	Brand         string    `gorm:"column:brand;type:varchar;not null" json:"brand"`
	LastFour      string    `gorm:"column:last_four;type:varchar(4);not null" json:"last_four"`
	ExpMonth      int16     `gorm:"column:exp_month;type:int;not null;check:exp_month>=1 AND exp_month<=12" json:"exp_month"`
	ExpYear       int16     `gorm:"column:exp_year;type:int;not null" json:"exp_year"`
	IsDefault     bool      `gorm:"column:is_default;type:boolean;not null;default:false" json:"is_default"`
}

func (Method) TableName() string {
	return "payment_methods"
}
