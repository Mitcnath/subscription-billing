package accounts

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamptz" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamptz" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at,omitempty"`
}

type UserAccounts struct {
	Base
	Email        string `gorm:"column:email;type:varchar;not null;unique" json:"email"`
	Username     string `gorm:"column:username;type:varchar;not null" json:"username"`
	PasswordHash string `gorm:"column:password_hash;type:varchar;not null" json:"-"`
}

func (UserAccounts) TableName() string {
	return "user_accounts"
}

// UserAccountReadModel is a read-only representation of the UserAccounts model, excluding sensitive fields like PasswordHash.
type UserAccountsReadModel struct {
	Base
	Email    string `gorm:"column:email;type:varchar;not null;unique" json:"email"`
	Username string `gorm:"column:username;type:varchar;not null" json:"username"`
}

func (UserAccountsReadModel) TableName() string {
	return "user_accounts"
}
