package services

import (
	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/repositories"
)

type ProductService struct{ Repo *repositories.ProductRepository }

func NewProductService(r *repositories.ProductRepository) *ProductService { return &ProductService{Repo: r} }

func (s *ProductService) List() ([]models.Product, error) { return s.Repo.List() }

func (s *ProductService) Create(p *models.Product) error { return s.Repo.Create(p) }

func (s *ProductService) Update(p *models.Product) error { return s.Repo.Update(p) }

func (s *ProductService) Delete(id uint) error { return s.Repo.Delete(id) }
