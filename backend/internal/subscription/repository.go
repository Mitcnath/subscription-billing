package subscription

import "gorm.io/gorm"

type Repository interface {
	FindReadModelByID(id int64) (*SubscriptionsReadModel, error)
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

// FindByID implements [Repository].
func (repository *repositoryImpl) FindReadModelByID(id int64) (*SubscriptionsReadModel, error) {
	var result SubscriptionsReadModel

	err := repository.db.
		Table("subscriptions s").
		Select(`
			s.id, s.created_at, s.updated_at,
			u.email AS user_account_email,
			p.name AS plan_name, p.amount AS plan_price_amount, p.currency AS plan_price_currency, p.billing_interval AS plan_billing_interval,
			s.trial_ends_at,
			s.current_period_ends_at,
			s.cancel_at_period_end,
			s.status,
			s.cancelled_at
		`).
		Joins("JOIN user_accounts u on u.id = s.user_account_id").
		Joins("JOIN subscription_plans p ON p.id = s.subscription_plan_id").
		Where("s.id = ?", id).
		Scan(&result).Error

	return &result, err
}
