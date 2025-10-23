package repositories

import (
	"errors"

	"mini-ecommerce/backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct { DB *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{DB: db} }

func (r *UserRepository) ByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(u *models.User) error { return r.DB.Create(u).Error }

func (r *UserRepository) EnsureAdmin(email, hash string) (*models.User, error) {
	var u models.User
	err := r.DB.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		u = models.User{Email: email, PasswordHash: hash, Role: models.RoleAdmin}
		if err := r.DB.Create(&u).Error; err != nil { return nil, err }
		return &u, nil
	}
	return &u, err
}
