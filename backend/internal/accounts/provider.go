package accounts

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Provider struct {
	accountsRepository Repository
}

func NewProvider(accountsRepository Repository) *Provider {
	return &Provider{accountsRepository: accountsRepository}
}

// TODO: Either remove this endpoint or add admin authorization to it

// GetAccounts godoc.
// @Summary      List user accounts
// @Tags         accounts
// @Produce      json
// @Param        limit    query     int     false  "Number of results per page (max 100)"  default(10)
// @Param        page     query     int     false  "Page number"                           default(1)
// @Param        sort_by  query     string  false  "Field to sort by"  Enums(username, email, created_at)  default(created_at)
// @Param        order    query     string  false  "Sort direction"    Enums(asc, desc)                    default(asc)
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/accounts [get]
func (provider *Provider) GetAccounts(ginContext *gin.Context) {
	limit, err := strconv.Atoi(ginContext.DefaultQuery("limit", "10"))

	if err != nil || limit <= 0 {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	if limit > 100 {
		limit = 100 // max limit to prevent abuse and performance issues
	}

	page, err := strconv.Atoi(ginContext.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}

	sortBy := ginContext.DefaultQuery("sort_by", "created_at")
	order := ginContext.DefaultQuery("order", "asc")

	results, total, err := provider.accountsRepository.FindAll(limit, page, sortBy, order)
	if err != nil {
		slog.Error("failed to fetch accounts", "error", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"data":       results,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

// Register godoc.
// @Summary      Register a new user account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Register request"
// @Success      201      {object}  UserAccounts
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/accounts/register [post]
func (provider *Provider) Register(ginContext *gin.Context) {
	var registerRequest RegisterRequest

	if err := ginContext.ShouldBindJSON(&registerRequest); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAccount, err := NewUserAccountService(provider.accountsRepository).Register(registerRequest)

	if err != nil {
		if errors.Is(err, EmailAlreadyTakenErr) {
			ginContext.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, UsernameAlreadyTakenErr) {
			ginContext.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": UnexpectedErr})
		return
	}

	ginContext.JSON(http.StatusCreated, userAccount)
}

// LoginRequest defines the request body for authenticating a user.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc.
// @Summary      Authenticate a user and return a JWT
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login request"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/accounts/login [post]
func (provider *Provider) Login(ginContext *gin.Context) {
	var loginRequest LoginRequest

	if err := ginContext.ShouldBindJSON(&loginRequest); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	signed, err := NewUserAccountService(provider.accountsRepository).Login(loginRequest.Email, loginRequest.Password)

	if err != nil {
		if errors.Is(err, EmailNotFoundErr) {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, InvalidEmailOrPasswordErr) {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"token": signed})
}

// GetMe godoc.
// @Summary      Get the authenticated user's account
// @Tags         accounts
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  UserAccounts
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/accounts/me [get]
func (provider *Provider) GetMe(ginContext *gin.Context) {
	userIDString := ginContext.GetString("userID")

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
		return
	}

	userAccount, err := provider.accountsRepository.FindReadModelByID(userID)
	if err != nil {
		slog.Error("failed to fetch user account", "error", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		return
	}

	ginContext.JSON(http.StatusOK, userAccount)
}
