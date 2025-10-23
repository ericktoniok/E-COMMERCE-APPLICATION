package repositories

import (
	"mini-ecommerce/backend/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository struct{ DB *gorm.DB }

func NewTransactionRepository(db *gorm.DB) *TransactionRepository { return &TransactionRepository{DB: db} }

func (r *TransactionRepository) Create(t *models.Transaction) error { return r.DB.Create(t).Error }

func (r *TransactionRepository) UpdateStatusByProviderRef(ref string, status models.TxStatus, raw string) error {
	return r.DB.Model(&models.Transaction{}).Where("provider_ref = ?", ref).Updates(map[string]interface{}{
		"status":     status,
		"raw_payload": raw,
	}).Error
}

func (r *TransactionRepository) ByProviderRef(ref string) (*models.Transaction, error) {
	var t models.Transaction
	if err := r.DB.Where("provider_ref = ?", ref).First(&t).Error; err != nil { return nil, err }
	return &t, nil
}
