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
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:200;not null"`
	Description string    `gorm:"type:text"`
	PriceCents  int64     `gorm:"not null"`
	Stock       int       `gorm:"not null;default:0"`
	ImageURL    string    `gorm:"size:500"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
