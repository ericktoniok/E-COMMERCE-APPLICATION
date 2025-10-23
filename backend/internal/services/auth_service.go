package services

import (
	"errors"
	"strings"

	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/repositories"
	"mini-ecommerce/backend/internal/utils"
)

type AuthService struct { Users *repositories.UserRepository }

func NewAuthService(users *repositories.UserRepository) *AuthService { return &AuthService{Users: users} }

func (s *AuthService) Register(email, password string) (*models.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || len(password) < 6 { return nil, errors.New("invalid payload") }
	h, err := utils.HashPassword(password)
	if err != nil { return nil, err }
	u := &models.User{Email: email, PasswordHash: h, Role: models.RoleCustomer}
	if err := s.Users.Create(u); err != nil { return nil, err }
	return u, nil
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	u, err := s.Users.ByEmail(strings.TrimSpace(strings.ToLower(email)))
	if err != nil { return nil, errors.New("invalid credentials") }
	if !utils.CheckPassword(u.PasswordHash, password) { return nil, errors.New("invalid credentials") }
	return u, nil
}
