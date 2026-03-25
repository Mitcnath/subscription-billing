package plans

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Provider struct is essentially a java class with only fields
type Provider struct {
	plansRepository PlansRepository
}

// Constructor
func NewProvider(plansRepository PlansRepository) *Provider {
	return &Provider{plansRepository: plansRepository}
}

// GetPlanById godoc
// @Summary      Get a subscription plan by ID
// @Tags         plans
// @Produce      json
// @Param        id   path      int  true  "Plan ID"
// @Success      200  {object}  map[string]SubscriptionPlans
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/plans/{id} [get]
func (provider *Provider) GetPlanById(ginContext *gin.Context) {
	id, err := strconv.ParseInt(ginContext.Param("id"), 10, 16)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid"})
		return
	}

	// Assignment and validation in the same statement
	// Checks if any error exists concerning the query
	result, err := provider.plansRepository.FindByID(int16(id))
	if err != nil {
		// Checks what the error actually is
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
			return
		}
		// Gives structured log if error is not caught
		slog.Error("failed to fetch plan", "error", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		return

	}
	ginContext.JSON(http.StatusOK, gin.H{"data": result})
}

// GetPlans godoc
// @Summary      List subscription plans
// @Tags         plans
// @Produce      json
// @Param        limit   query     int     false  "Number of results per page"  default(10)
// @Param        after   query     int     false  "Cursor: last seen plan ID for pagination"
// @Param        status  query     string  false  "Filter by plan status"  Enums(active, deprecated)  default(active)
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api/v1/plans [get]
func (provider *Provider) GetPlans(ginContext *gin.Context) {

	// Paginate data going out so the entire table isn't fetched
	// Parse limit using string convert ASCII to integer standard library
	limit, err := strconv.Atoi(ginContext.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	// Using '_' to discard 'err' since because I don't intend to use it
	// strconv.ParseInt(arg, use base10 numbering, 16-bit Integer to match the model struct)
	after, _ := strconv.ParseInt(ginContext.DefaultQuery("after", "0"), 10, 16)

	// Cast PlanStatus type so validator can be used after
	planStatus := PlanStatus(ginContext.DefaultQuery("status", "active"))

	if !planStatus.Valid() {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	// Creates an empty, initialised slice of SubscriptionPlans with a length of 0
	// so 'null' will never be returned in the JSON response
	// https://go.dev/tour/moretypes/13
	results, err := provider.plansRepository.FindAll(limit, int16(after), planStatus)

	// https://gorm.io/docs/query.html#Order
	if err != nil {
		slog.Error("failed to fetch plans", "error", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		return
	}

	// Not a true cursor, just a naming convention. But its worth looking into for large data sets.
	// https://www.postgresql.org/docs/current/plpgsql-cursors.html
	var nextCursor *int16
	if len(results) == limit {
		last := results[len(results)-1].ID
		nextCursor = &last
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"data": results,
		// Next 'after' position. Will be null on final page
		"next_cursor": nextCursor,
	})
}

// CreatePlanRequest defines the request body for creating a subscription plan.
type CreatePlanRequest struct {
	Name            string          `json:"name" binding:"required"`
	Description     string          `json:"description" binding:"required"`
	Amount          uint64          `json:"amount" binding:"required,gt=0"`
	Currency        string          `json:"currency" binding:"required,iso4217"`
	BillingInterval BillingInterval `json:"billing_interval" binding:"required"`
}

// CreatePlan godoc
// @Summary      Create a subscription plan
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        request  body      CreatePlanRequest  true  "Create plan request"
// @Success      201      {object}  SubscriptionPlans
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/plans/create [post]
func (provider *Provider) CreatePlan(ginContext *gin.Context) {

	var subscriptionPlanRequest CreatePlanRequest

	// ShouldBindJSON allows the error to be handled manually
	if err := ginContext.ShouldBindJSON(&subscriptionPlanRequest); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !subscriptionPlanRequest.BillingInterval.Valid() {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "billing interval does not exist. Must be one of the following: \"daily\", \"weekly\", \"bi_weekly\", \"monthly\", \"quarterly\", \"semi_annual\", \"annual\""})
		return
	}

	plan := SubscriptionPlans{
		Name:            subscriptionPlanRequest.Name,
		Description:     subscriptionPlanRequest.Description,
		Amount:          subscriptionPlanRequest.Amount,
		Currency:        subscriptionPlanRequest.Currency,
		BillingInterval: subscriptionPlanRequest.BillingInterval,
	}

	if err := provider.plansRepository.Create(&plan); err != nil {
		// Name has a unique constraint at db level so check that using gorm
		// https://pkg.go.dev/errors
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ginContext.JSON(http.StatusConflict, gin.H{"error": "a plan with this name already exists"})
			return
		}
		// Structured logging: https://pkg.go.dev/log/slog
		slog.Error("failed to create plan", "error", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		return
	}

	ginContext.JSON(http.StatusCreated, plan)
}

// UpdatePlanRequest defines the request body for updating a subscription plan.
type UpdatePlanRequest struct {
	Name            *string          `json:"name"`
	Description     *string          `json:"description"`
	Amount          *uint64          `json:"amount" binding:"omitempty,gt=0"`
	Currency        *string          `json:"currency" binding:"omitempty,iso4217"`
	BillingInterval *BillingInterval `json:"billing_interval"`
}

// UpdatePlanById godoc
// @Summary      Update a subscription plan
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "Plan ID"
// @Param        request  body      UpdatePlanRequest  true  "Update plan request — all fields optional"
// @Success      200      {object}  SubscriptionPlans
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/plans/update/{id} [patch]
func (provider *Provider) UpdatePlanById(ginContext *gin.Context) {

	id, err := strconv.ParseInt(ginContext.Param("id"), 10, 16)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid"})
		return
	}

	// Use pointers in the request struct. A nil pointer means the field isn't provided and will be left as-is. This is necessary if the fields are not required
	var subscriptionPlanRequest UpdatePlanRequest

	if err := ginContext.ShouldBindJSON(&subscriptionPlanRequest); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if (subscriptionPlanRequest == UpdatePlanRequest{}) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "at least one field must be provided to update"})
		return
	}

	if subscriptionPlanRequest.BillingInterval != nil && !subscriptionPlanRequest.BillingInterval.Valid() {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "billing interval does not exist. Must be one of the following: \"daily\", \"weekly\", \"bi_weekly\", \"monthly\", \"quarterly\", \"semi_annual\", \"annual\""})
		return
	}

	result, err := provider.plansRepository.FindByID(int16(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
			return
		}
		slog.Error("failed to fetch plan", "error", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		return
	}

	if result.Status == PlanStatusDeprecated {
		ginContext.JSON(http.StatusConflict, gin.H{"error": "cannot update a deprecated plan. Please change the plan status to active before updating"})
		return
	}

	if err := provider.plansRepository.Update(result, subscriptionPlanRequest); err != nil {
		slog.Error("failed to update plan", "error", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		return
	}

	ginContext.JSON(http.StatusOK, result)
}

// UpdatePlanStatusRequest defines the request body for updating a plan's status.
type UpdatePlanStatusRequest struct {
	Status PlanStatus `json:"status"`
}

// UpdatePlanStatusById godoc
// @Summary      Update a plan's status
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id       path      int                      true  "Plan ID"
// @Param        request  body      UpdatePlanStatusRequest  true  "Update plan status request"
// @Success      200      {object}  SubscriptionPlans
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/plans/update/status/{id} [patch]
func (provider *Provider) UpdatePlanStatusById(ginContext *gin.Context) {

	id, err := strconv.ParseInt(ginContext.Param("id"), 10, 16)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid"})
		return
	}

	var subscriptionPlanRequest UpdatePlanStatusRequest

	if err := ginContext.ShouldBindJSON(&subscriptionPlanRequest); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !subscriptionPlanRequest.Status.Valid() {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "status does not exist. Must be one of the following: \"active\", \"deprecated\""})
		return
	}

	result, err := provider.plansRepository.FindByID(int16(id))

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
			return
		}
		slog.Error("failed to fetch plan", "error", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		return
	}

	if err := provider.plansRepository.UpdateStatus(result, subscriptionPlanRequest); err != nil {
		slog.Error("failed to update plan", "error", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		return
	}

	ginContext.JSON(http.StatusOK, result)
}
