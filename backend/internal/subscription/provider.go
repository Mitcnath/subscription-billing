package subscription

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Provider struct {
	subscriptionRepository Repository
}

func NewProvider(subscriptionRepository Repository) *Provider {
	return &Provider{subscriptionRepository: subscriptionRepository}
}

// GetSubscriptionByID godoc.
//
//	@Summary	Get a subscription by ID
//	@Tags		subscriptions
//	@Produce	json
//	@Param		id	path		int								true	"Subscription ID"
//	@Success	200	{object}	subscription.SubscriptionsReadModel
//	@Failure	400	{object}	map[string]string
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Security	BearerAuth
//	@Router		/api/v1/subscriptions/{id} [get]
func (provider *Provider) GetSubscriptionByID(ginContext *gin.Context) {
	id, err := strconv.ParseInt(ginContext.Param("id"), 10, 64)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid"})
		return
	}

	result, err := provider.subscriptionRepository.FindReadModelByID(id)

	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, result)
}
