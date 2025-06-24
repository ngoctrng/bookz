package repository

import (
	"github.com/ngoctrng/bookz/internal/account"
	"gorm.io/gorm"
)

const tblusers = "users"

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(user *account.User) error {
	schema := DomainToUserSchema(user)
	if err := r.db.Table(tblusers).Create(schema).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindByEmail(email string) (*account.User, error) {
	var u UserSchema

	err := r.db.Table(tblusers).Where("email = ?", email).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return UserSchemaToDomain(&u), nil
}
