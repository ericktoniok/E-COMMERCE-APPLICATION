package models

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

type OrderStatus string

const (
	OrderPending OrderStatus = "PENDING"
	OrderPaid    OrderStatus = "PAID"
	OrderFailed  OrderStatus = "FAILED"
)

type TxStatus string

const (
	TxPending TxStatus = "PENDING"
	TxSuccess TxStatus = "SUCCESS"
	TxFailed  TxStatus = "FAILED"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	Role         Role      `gorm:"size:20;not null;default:customer"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	PriceCents  int64     `gorm:"not null" json:"price_cents"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	ImageURL    string    `gorm:"size:500" json:"image_url"`
	Category    string    `gorm:"size:120" json:"category"`
	SKU         string    `gorm:"size:64" json:"sku"`
	Rating      float32   `gorm:"type:decimal(3,2);default:0" json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Order struct {
	ID         uint        `gorm:"primaryKey"`
	UserID     uint        `gorm:"index;not null"`
	User       User        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status     OrderStatus `gorm:"size:16;not null;default:PENDING"`
	TotalCents int64       `gorm:"not null;default:0"`
	Items      []OrderItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrderItem struct {
	ID         uint    `gorm:"primaryKey"`
	OrderID    uint    `gorm:"index;not null"`
	ProductID  uint    `gorm:"index;not null"`
	Product    Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Qty        int     `gorm:"not null"`
	PriceCents int64   `gorm:"not null"`
}

type Transaction struct {
	ID           uint      `gorm:"primaryKey"`
	OrderID      uint      `gorm:"index;not null"`
	Order        Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProviderRef  string    `gorm:"size:100;index"`
	Status       TxStatus  `gorm:"size:16;not null;default:PENDING"`
	AmountCents  int64     `gorm:"not null"`
	RawPayload   string    `gorm:"type:text"`
	CreatedAt    time.Time
}
