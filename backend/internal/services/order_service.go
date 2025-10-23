package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/repositories"
)

type CheckoutItem struct {
	ProductID uint `json:"product_id"`
	Qty       int  `json:"qty"`
}

type MpesaClient struct {
	BaseURL        string
	BackendWebhook string
	HTTP           *http.Client
}

type OrderService struct {
	Products    *repositories.ProductRepository
	Orders      *repositories.OrderRepository
	Transactions *repositories.TransactionRepository
	Mpesa       *MpesaClient
}

func NewOrderService(p *repositories.ProductRepository, o *repositories.OrderRepository, t *repositories.TransactionRepository, m *MpesaClient) *OrderService {
	return &OrderService{Products: p, Orders: o, Transactions: t, Mpesa: m}
}

type mpesaPayReq struct {
	OrderID     string `json:"order_id"`
	AmountCents int64  `json:"amount_cents"`
}

type mpesaPayResp struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

func (s *OrderService) Checkout(userID uint, items []CheckoutItem) (*models.Order, *models.Transaction, error) {
	// compute total and build order items
	var orderItems []models.OrderItem
	var total int64
	for _, it := range items {
		p, err := s.Products.Get(it.ProductID)
		if err != nil { return nil, nil, err }
		line := models.OrderItem{ProductID: p.ID, Qty: it.Qty, PriceCents: p.PriceCents}
		orderItems = append(orderItems, line)
		total += int64(it.Qty) * p.PriceCents
	}
	o := &models.Order{UserID: userID, Status: models.OrderPending, TotalCents: total}
	if err := s.Orders.CreateWithItems(o, orderItems); err != nil { return nil, nil, err }

	// call mock mpesa
	payload := mpesaPayReq{OrderID: fmt.Sprintf("%d", o.ID), AmountCents: total}
	buf, _ := json.Marshal(payload)
	url := s.Mpesa.BaseURL + "/pay"
	resp, err := s.Mpesa.HTTP.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil { return o, nil, err }
	defer resp.Body.Close()
	var pr mpesaPayResp
	_ = json.NewDecoder(resp.Body).Decode(&pr)

	// store transaction
	tx := &models.Transaction{OrderID: o.ID, ProviderRef: pr.TransactionID, Status: models.TxPending, AmountCents: total}
	if err := s.Transactions.Create(tx); err != nil { return o, nil, err }
	return o, tx, nil
}

func (s *OrderService) ListByUser(userID uint) ([]models.Order, error) {
	return s.Orders.ByUser(userID)
}

func (s *OrderService) ListAll() ([]models.Order, error) {
	return s.Orders.All()
}
