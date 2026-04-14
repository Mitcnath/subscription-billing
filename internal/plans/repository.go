package plans

import (
	"time"

	"gorm.io/gorm"
)

// https://medium.com/@joshuasajeevnv/how-to-implement-repository-pattern-in-golang-go-ff2625fe407f
type Repository interface {
	FindByID(id int16) (*SubscriptionPlans, error)
	FindAll(limit int, after int16, status PlanStatus) ([]SubscriptionPlans, error)
	FindLast() (SubscriptionPlans, error)
	Create(SubscriptionPlan *SubscriptionPlans) error
	Update(SubscriptionPlan *SubscriptionPlans, req UpdatePlanRequest) error
	DeprecatedPlan(SubscriptionPlan *SubscriptionPlans) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

// Create implements [Repository].
func (repository *repositoryImpl) Create(subscriptionPlan *SubscriptionPlans) error {
	return repository.db.Create(subscriptionPlan).Error
}

// FindAll implements [Repository].
func (repository *repositoryImpl) FindAll(limit int, after int16, status PlanStatus) ([]SubscriptionPlans, error) {
	results := make([]SubscriptionPlans, 0)
	query := repository.db.Order("id ASC").Limit(limit).Where("status = ?", status)

	if after != 0 {
		query = query.Where("id > ?", after)
	}

	err := query.Find(&results).Error
	return results, err
}

func (repository *repositoryImpl) FindLast() (SubscriptionPlans, error) {
	var result SubscriptionPlans
	err := repository.db.Last(&result).Error
	return result, err
}

// FindByID implements [Repository].
func (repository *repositoryImpl) FindByID(id int16) (*SubscriptionPlans, error) {
	var result SubscriptionPlans
	err := repository.db.First(&result, id).Error
	return &result, err
}

// Update implements [Repository].
func (repository *repositoryImpl) Update(subscriptionPlan *SubscriptionPlans, req UpdatePlanRequest) error {
	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.BillingInterval != nil {
		updates["billing_interval"] = *req.BillingInterval
	}
	if req.Money != nil {
		updates["amount"] = req.Money.Amount
		updates["currency"] = req.Money.Currency
	}

	subscriptionPlan.UpdatedAt = time.Now()
	return repository.db.Model(subscriptionPlan).Updates(updates).Error
}

func (repository *repositoryImpl) DeprecatedPlan(subscriptionPlan *SubscriptionPlans) error {
	if err := subscriptionPlan.Deprecate(); err != nil {
		return err
	}
	subscriptionPlan.UpdatedAt = time.Now()
	return repository.db.Model(subscriptionPlan).Update("status", PlanStatusDeprecated).Error
}
