package accounts

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAccountService struct {
	userAccounts Repository
}

func NewUserAccountService(userAccounts Repository) *UserAccountService {
	return &UserAccountService{userAccounts: userAccounts}
}

var (
	EmailAlreadyTakenErr    = errors.New("email already in use")
	UsernameAlreadyTakenErr = errors.New("username already in use")
	UnexpectedErr           = errors.New("an unexpected error has occured")
	EmailNotFoundErr        = errors.New("email not found")
)

// RegisterRequest defines the request body for creating a user account.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"` // Minimum 8 characters for basic password strength
}

func (service *UserAccountService) Register(request RegisterRequest) (*UserAccountsReadModel, error) {

	_, err := service.userAccounts.FindReadModelByEmail(request.Email)

	if err == nil {
		slog.Warn("registration attempt with existing email", "email", request.Email)
		return nil, EmailAlreadyTakenErr
	}

	_, err = service.userAccounts.FindReadModelByUsername(request.Username)

	if err == nil {
		slog.Warn("registraton attempt with existing username", "username", request.Username)
		return nil, UsernameAlreadyTakenErr
	}

	hashedPassword, err := hashPassword(request.Password)

	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, UnexpectedErr
	}

	userAccount := UserAccounts{
		Email:        request.Email,
		Username:     request.Username,
		PasswordHash: hashedPassword,
	}

	if err := service.userAccounts.Create(&userAccount); err != nil {
		slog.Error("failed to create user account", "error", err)
		return nil, UnexpectedErr
	}

	userAccountReadModel, err := service.userAccounts.FindReadModelByEmail(request.Email)

	if err != nil {
		slog.Error("account creation error", "account", err)
		return nil, UnexpectedErr
	}

	return userAccountReadModel, nil
}

var (
	InvalidPasswordErr        = errors.New("failed to verify password")
	InvalidEmailOrPasswordErr = errors.New("invalid email or password")
	JWTSigningErr             = errors.New("failed to sign token")
)

func (service *UserAccountService) Login(email string, password string) (string, error) {

	userAccount, err := service.userAccounts.FindByEmail(email)

	if err != nil {
		slog.Error("email lookup failed", "error", err)
		return "", EmailNotFoundErr
	}

	match, err := comparePassword(password, userAccount.PasswordHash)

	if err != nil {
		slog.Error("comparison function aborted", "error", err)
		return "", InvalidEmailOrPasswordErr
	}

	if !match {
		slog.Warn("password match failed", "password", InvalidPasswordErr)
		return "", InvalidEmailOrPasswordErr
	}

	// Create JWT token
	claims := jwt.RegisteredClaims{
		Subject:   userAccount.ID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // Token valid for 1 hour
	}

	// Create the token using the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		slog.Error("failed to sign token", "error", err)
		return "", UnexpectedErr
	}

	return signed, nil
}

func (service *UserAccountService) GetAccounts(limit int, page int, sortBy string, order string) ([]UserAccountsReadModel, int64, error) {

	// Validate and sanitize sortBy and order parameters to prevent SQL injection
	allowedFields := map[string]string{
		"username":   "username",
		"email":      "email",
		"created_at": "created_at",
	}

	// Default to "created_at" if sortBy is not valid
	column, ok := allowedFields[sortBy]
	if !ok {
		column = "created_at"
	}

	// Default to "asc" if order is not valid
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	return service.userAccounts.FindAll(limit, page, column, order)
}

type UpdateUserAccountRequest struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
}

func (service *UserAccountService) UpdateAccountByID(id uuid.UUID, userAccountRequest UpdateUserAccountRequest) (*UserAccountsReadModel, error) {

	account, err := service.userAccounts.FindReadModelByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("record not found", "error", gorm.ErrRecordNotFound)
			return nil, fmt.Errorf("user account not found")
		}
	}

	if userAccountRequest.Email != nil {
		account.Email = *userAccountRequest.Email
	}

	if userAccountRequest.Username != nil {
		account.Username = *userAccountRequest.Username
	}

	if err := service.userAccounts.UpdateReadModel(account); err != nil {
		return nil, err
	}

	account.UpdatedAt = time.Now()
	return service.userAccounts.FindReadModelByID(account.ID)
}
