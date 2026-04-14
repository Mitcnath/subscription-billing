package plans

import (
	"billingService/backend/internal/money"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type PlansService struct {
	plans Repository
}

func NewPlansService(plans Repository) *PlansService {
	return &PlansService{plans: plans}
}

var (
	InvalidMoneyErr           = errors.New("invalid money: amount must be non-negative and currency must be a valid ISO 4217 code")
	InvalidBillingIntervalErr = errors.New("billing interval does not exist. Must be one of the following: \"daily\", \"weekly\", \"bi_weekly\", \"monthly\", \"quarterly\", \"semi_annual\", \"annual\"")
	ErrDuplicatedKey          = errors.New("a plan with this name already exists")
	ErrCreationFailed         = errors.New("failed to create plan")
	ErrPlanNotFound           = errors.New("plan not found")
)

// CreatePlanRequest defines the request body for creating a subscription plan.
type CreatePlanRequest struct {
	Name            string          `json:"name" binding:"required"`
	Description     string          `json:"description" binding:"required"`
	Money           money.Money     `json:"money" binding:"required"`
	BillingInterval BillingInterval `json:"billing_interval" binding:"required"`
}

func (service *PlansService) CreatePlan(request CreatePlanRequest) (*SubscriptionPlans, error) {
	if !request.Money.Valid() {
		return nil, InvalidMoneyErr
	}
	if !request.BillingInterval.Valid() {
		return nil, InvalidBillingIntervalErr
	}

	planToCreate := SubscriptionPlans{
		Name:            request.Name,
		Description:     request.Description,
		Price:           request.Money,
		BillingInterval: request.BillingInterval,
	}

	if err := service.plans.Create(&planToCreate); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			slog.Warn("creation attempt with existing plan", "plan", err)
			return nil, ErrDuplicatedKey
		}
		slog.Error("failed to create plan", "error", ErrCreationFailed)
		return nil, fmt.Errorf("an unexpected error has occurred")
	}

	plan, err := service.plans.FindLast()

	if err != nil {
		return nil, ErrPlanNotFound
	}

	return &plan, nil
}

var (
	ErrNoChangesMade      = errors.New("at least one field must be provided to update")
	ErrCannotBeDeprecated = errors.New("cannot update a deprecated plan. Please change the plan status to active before updating")
)

func (service *PlansService) UpdatePlan(request UpdatePlanRequest, id int16) (*SubscriptionPlans, error) {

	if (request == UpdatePlanRequest{}) {
		return nil, ErrNoChangesMade
	}
	if request.Money != nil && !request.Money.Valid() {
		return nil, InvalidMoneyErr
	}
	if request.BillingInterval != nil && !request.BillingInterval.Valid() {
		return nil, InvalidBillingIntervalErr
	}

	existingPlan, err := service.plans.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("plan not found", "plan", err)
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("an unexpected error has occurred")
	}

	if !existingPlan.CanBeDeprecated() {
		slog.Warn("deprecation attempt on already deprecated plan", "status", ErrCannotBeDeprecated)
		return nil, ErrCannotBeDeprecated
	}

	if err := service.plans.Update(existingPlan, request); err != nil {
		slog.Error("failed to update plan", "error", err)
		return nil, fmt.Errorf("an unexpected error has occurred")
	}

	return service.plans.FindByID(id)

}

func (service *PlansService) DeprecatePlanByID(id int16) (*SubscriptionPlans, error) {
	existingPlan, err := service.plans.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("plan not found", "plan", err)
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("an unexpected error has occurred")
	}
	if !existingPlan.CanBeDeprecated() {
		slog.Warn("deprecation attempt on already deprecated plan", "status", ErrCannotBeDeprecated)
		return nil, ErrCannotBeDeprecated
	}
	if err := service.plans.DeprecatedPlan(existingPlan); err != nil {
		slog.Error("failed to update plan", "error", err)
		return nil, fmt.Errorf("an unexpected error has occurred")
	}
	return service.plans.FindByID(id)
}
