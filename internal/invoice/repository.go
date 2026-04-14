package invoice

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(invoice *Invoice) error
	FindByID(id int64) (*Invoice, error)
	FindReadModelByID(id int64) (*InvoiceReadModel, error)
	Update(invoice *Invoice) error
}

type repositoryImpl struct {
	db *gorm.DB
}

// Create implements [Repository].
func (repository *repositoryImpl) Create(invoice *Invoice) error {
	return repository.db.Create(invoice).Error
}

func (repository *repositoryImpl) FindByID(id int64) (*Invoice, error) {
	var result Invoice
	err := repository.db.First(&result, id).Error
	return &result, err
}

// FindReadModelByID implements [Repository].
func (repository *repositoryImpl) FindReadModelByID(id int64) (*InvoiceReadModel, error) {
	var result InvoiceReadModel

	// https://gorm.io/docs/query.html#Joins
	err := repository.db.
		Table("invoices i").
		Select(`
			i.id, i.created_at, i.updated_at,
			u.email AS user_email,
			sp.name AS subscription_name,
			sp.amount AS subscription_price_amount, sp.currency AS subscription_price_currency,
			i.status,
			i.amount AS paid_amount, i.currency AS paid_currency,
			i.pdf_url
		`).
		Joins("JOIN user_accounts u ON u.id = i.user_account_id").
		Joins("JOIN subscriptions s ON s.id = i.subscription_id").
		Joins("JOIN subscription_plans sp ON sp.id = s.subscription_plan_id").
		Where("i.id = ?", id).
		Scan(&result).Error

	return &result, err
}

// Update implements [Repository].
func (repository *repositoryImpl) Update(invoice *Invoice) error {
	return errors.New("Updated not implemented")
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}
