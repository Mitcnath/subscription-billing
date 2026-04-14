package accounts

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindReadModelByID(id uuid.UUID) (*UserAccountsReadModel, error)
	FindByID(id uuid.UUID) (*UserAccounts, error)
	FindReadModelByUsername(username string) (*UserAccountsReadModel, error)
	FindReadModelByEmail(email string) (*UserAccountsReadModel, error)
	FindByEmail(email string) (*UserAccounts, error)
	FindAll(limit int, page int, sortBy string, order string) ([]UserAccountsReadModel, int64, error)
	Create(account *UserAccounts) error
	UpdateReadModel(account *UserAccountsReadModel) error
	Update(account *UserAccounts) error
}

type accountsRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountsRepository(db *gorm.DB) Repository {
	return &accountsRepositoryImpl{db: db}
}

func (repository *accountsRepositoryImpl) FindReadModelByEmail(email string) (*UserAccountsReadModel, error) {
	var result UserAccountsReadModel
	err := repository.db.Where("email = ?", email).First(&result).Error
	return &result, err
}

func (repository *accountsRepositoryImpl) FindByEmail(email string) (*UserAccounts, error) {
	var result UserAccounts
	err := repository.db.Where("email = ?", email).First(&result).Error
	return &result, err
}

func (repository *accountsRepositoryImpl) FindReadModelByID(id uuid.UUID) (*UserAccountsReadModel, error) {
	var result UserAccountsReadModel
	err := repository.db.First(&result, "id = ?", id).Error
	return &result, err
}

func (repository *accountsRepositoryImpl) FindByID(id uuid.UUID) (*UserAccounts, error) {
	var result UserAccounts
	err := repository.db.First(&result, "id = ?", id).Error
	return &result, err
}

func (repository *accountsRepositoryImpl) FindReadModelByUsername(username string) (*UserAccountsReadModel, error) {
	var result UserAccountsReadModel
	err := repository.db.Where("username = ?", username).First(&result).Error
	return &result, err
}

func (repository *accountsRepositoryImpl) Create(account *UserAccounts) error {
	return repository.db.Create(account).Error
}

func (repository *accountsRepositoryImpl) UpdateReadModel(account *UserAccountsReadModel) error {
	return repository.db.Where("id = ?", account.ID).Updates(account).Error
}

func (repository *accountsRepositoryImpl) Update(account *UserAccounts) error {
	return repository.db.Where("id = ?", account.ID).Updates(account).Error
}

func (repository *accountsRepositoryImpl) FindAll(limit int, page int, sortBy string, order string) ([]UserAccountsReadModel, int64, error) {
	var results []UserAccountsReadModel
	var total int64

	// Default to "created_at"
	column := "created_at"

	// Calculate the offset for pagination based on the page number and limit
	offset := (page - 1) * limit

	// Get the total count of records for pagination
	if err := repository.db.Model(&UserAccountsReadModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Execute the query with sorting, pagination, and return the results
	err := repository.db.Order(column + " " + order).Limit(limit).Offset(offset).Find(&results).Error
	return results, total, err
}
