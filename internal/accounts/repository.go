package accounts

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountsRepository interface {
	FindByID(id uuid.UUID) (*UserAccounts, error)
	FindByUsername(username string) (*UserAccounts, error)
	FindByEmail(email string) (*UserAccounts, error)
	FindAll(limit int, page int, sortBy string, order string) ([]UserAccounts, int64, error)
	Create(account *UserAccounts) error
}

type accountsRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountsRepository(db *gorm.DB) AccountsRepository {
	return &accountsRepositoryImpl{db: db}
}

// FindByEmail implements [AccountsRepository].
func (repository *accountsRepositoryImpl) FindByEmail(email string) (*UserAccounts, error) {
	var result UserAccounts
	err := repository.db.Where("email = ?", email).First(&result).Error
	return &result, err
}

// FindByID implements [AccountsRepository].
func (repository *accountsRepositoryImpl) FindByID(id uuid.UUID) (*UserAccounts, error) {
	var result UserAccounts
	err := repository.db.First(&result, id).Error
	return &result, err
}

// FindByUsername implements [AccountsRepository].
func (repository *accountsRepositoryImpl) FindByUsername(username string) (*UserAccounts, error) {
	var result UserAccounts
	err := repository.db.Where("username = ?", username).First(&result).Error
	return &result, err
}

// Create implements [AccountsRepository].
func (repository *accountsRepositoryImpl) Create(account *UserAccounts) error {
	return repository.db.Create(account).Error
}

// FindAll implements [AccountsRepository].
func (repository *accountsRepositoryImpl) FindAll(limit int, page int, sortBy string, order string) ([]UserAccounts, int64, error) {
	var results []UserAccounts
	var total int64

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

	// Calculate the offset for pagination based on the page number and limit
	offset := (page - 1) * limit

	// Get the total count of records for pagination
	if err := repository.db.Model(&UserAccounts{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Execute the query with sorting, pagination, and return the results
	err := repository.db.Order(column + " " + order).Limit(limit).Offset(offset).Find(&results).Error
	return results, total, err
}
