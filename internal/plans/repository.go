package plans

import "gorm.io/gorm"

// https://medium.com/@joshuasajeevnv/how-to-implement-repository-pattern-in-golang-go-ff2625fe407f
type PlansRepository interface {
	FindByID(id int16) (*SubscriptionPlans, error)
	FindAll(limit int, after int16, status PlanStatus) ([]SubscriptionPlans, error)
	Create(SubscriptionPlan *SubscriptionPlans) error
	Update(SubscriptionPlan *SubscriptionPlans, req UpdatePlanRequest) error
	UpdateStatus(SubscriptionPlan *SubscriptionPlans, req UpdatePlanStatusRequest) error
}

type plansRepositoryImpl struct {
	db *gorm.DB
}

func NewPlansRepository(db *gorm.DB) PlansRepository {
	return &plansRepositoryImpl{db: db}
}

// Create implements [PlansRepository].
func (repository *plansRepositoryImpl) Create(SubscriptionPlan *SubscriptionPlans) error {
	return repository.db.Create(SubscriptionPlan).Error
}

// FindAll implements [PlansRepository].
func (repository *plansRepositoryImpl) FindAll(limit int, after int16, status PlanStatus) ([]SubscriptionPlans, error) {
	results := make([]SubscriptionPlans, 0)
	query := repository.db.Order("id ASC").Limit(limit).Where("status = ?", status)

	if after != 0 {
		query = query.Where("id > ?", after)
	}

	err := query.Find(&results).Error
	return results, err
}

// FindByID implements [PlansRepository].
func (repository *plansRepositoryImpl) FindByID(id int16) (*SubscriptionPlans, error) {
	var result SubscriptionPlans
	err := repository.db.First(&result, id).Error
	return &result, err
}

// Update implements [PlansRepository].
func (repository *plansRepositoryImpl) Update(SubscriptionPlan *SubscriptionPlans, req UpdatePlanRequest) error {
	return repository.db.Model(SubscriptionPlan).Updates(req).Error
}

// UpdateStatus implements [PlansRepository].
func (repository *plansRepositoryImpl) UpdateStatus(SubscriptionPlan *SubscriptionPlans, req UpdatePlanStatusRequest) error {
	return repository.db.Model(SubscriptionPlan).Updates(req).Error
}
