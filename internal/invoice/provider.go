package invoice

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Provider struct {
	invoiceRepository Repository
}

func NewProvider(invoiceRepository Repository) *Provider {
	return &Provider{invoiceRepository: invoiceRepository}
}

// GetInvoiceByID godoc.
//
//	@Summary	Get an invoice by ID
//	@Tags		invoices
//	@Produce	json
//	@Param		id	path		int						true	"Invoice ID"
//	@Success	200	{object}	invoice.InvoiceReadModel
//	@Failure	400	{object}	map[string]string
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Security	BearerAuth
//	@Router		/api/v1/invoices/{id} [get]
func (provider *Provider) GetInvoiceByID(ginContext *gin.Context) {
	id, err := strconv.ParseInt(ginContext.Param("id"), 10, 64)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid"})
		return
	}

	result, err := provider.invoiceRepository.FindReadModelByID(id)

	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, result)
}
