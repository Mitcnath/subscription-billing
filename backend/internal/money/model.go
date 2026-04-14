package money

import "github.com/rmg/iso4217"

type Money struct {
	Amount   int64  `gorm:"column:amount;type:bigint;check:amount>=0;not null" json:"amount"` // Amount in cents
	Currency string `gorm:"column:currency;type:varchar;not null" json:"currency"`
}

func (money Money) Valid() bool {
	if money.Amount < 0 {
		return false
	}
	if money.Currency == "" {
		return false
	}
	if code, _ := iso4217.ByName(money.Currency); code == 0 || code == 999 {
		return false
	}
	return true
}
