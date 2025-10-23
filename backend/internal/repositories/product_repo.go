package repositories

import (
	"mini-ecommerce/backend/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct{ DB *gorm.DB }

func NewProductRepository(db *gorm.DB) *ProductRepository { return &ProductRepository{DB: db} }

func (r *ProductRepository) List() ([]models.Product, error) {
	var items []models.Product
	if err := r.DB.Order("id desc").Find(&items).Error; err != nil { return nil, err }
	return items, nil
}

func (r *ProductRepository) Get(id uint) (*models.Product, error) {
	var p models.Product
	if err := r.DB.First(&p, id).Error; err != nil { return nil, err }
	return &p, nil
}

func (r *ProductRepository) Create(p *models.Product) error { return r.DB.Create(p).Error }
func (r *ProductRepository) Update(p *models.Product) error { return r.DB.Save(p).Error }
func (r *ProductRepository) Delete(id uint) error { return r.DB.Delete(&models.Product{}, id).Error }
