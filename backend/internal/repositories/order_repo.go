package repositories

import (
	"mini-ecommerce/backend/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct{ DB *gorm.DB }

func NewOrderRepository(db *gorm.DB) *OrderRepository { return &OrderRepository{DB: db} }

func (r *OrderRepository) CreateWithItems(o *models.Order, items []models.OrderItem) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil { return err }
		for i := range items {
			items[i].OrderID = o.ID
		}
		if err := tx.Create(&items).Error; err != nil { return err }
		return nil
	})
}

func (r *OrderRepository) ByUser(userID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := r.DB.Preload("Items").Order("id desc").Where("user_id = ?", userID).Find(&orders).Error; err != nil { return nil, err }
	return orders, nil
}

func (r *OrderRepository) All() ([]models.Order, error) {
	var orders []models.Order
	if err := r.DB.Preload("Items").Order("id desc").Find(&orders).Error; err != nil { return nil, err }
	return orders, nil
}

func (r *OrderRepository) UpdateStatus(id uint, status models.OrderStatus) error {
	return r.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) Get(id uint) (*models.Order, error) {
	var o models.Order
	if err := r.DB.Preload("Items").First(&o, id).Error; err != nil { return nil, err }
	return &o, nil
}
