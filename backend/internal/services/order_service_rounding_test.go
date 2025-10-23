package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/repositories"
	"mini-ecommerce/backend/internal/testutil"
)

func TestCheckout_RoundingAndLargeTotals(t *testing.T) {
	db := testutil.OpenTestDB(t)
	migrate(t, db)

	prodRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	// Seed products with mixed prices
	p1 := &models.Product{Name: "High", PriceCents: 19999, Stock: 10}   // $199.99
	p2 := &models.Product{Name: "Low", PriceCents: 99, Stock: 10}       // $0.99
	p3 := &models.Product{Name: "Mid", PriceCents: 12345, Stock: 10}    // $123.45
	if err := prodRepo.Create(p1); err != nil { t.Fatalf("seed p1: %v", err) }
	if err := prodRepo.Create(p2); err != nil { t.Fatalf("seed p2: %v", err) }
	if err := prodRepo.Create(p3); err != nil { t.Fatalf("seed p3: %v", err) }

	// Mock Mpesa
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"transaction_id": "tx-round", "status": "QUEUED"})
	}))
	defer ts.Close()

	svc := NewOrderService(prodRepo, orderRepo, txRepo, &MpesaClient{BaseURL: ts.URL, HTTP: ts.Client()})
	order, _, err := svc.Checkout(1, []CheckoutItem{{ProductID: p1.ID, Qty: 2}, {ProductID: p2.ID, Qty: 5}, {ProductID: p3.ID, Qty: 1}})
	if err != nil { t.Fatalf("checkout: %v", err) }
	// Expected total: 2*19999 + 5*99 + 1*12345 = 39998 + 495 + 12345 = 52838
	if want := int64(52838); order.TotalCents != want { t.Fatalf("total cents want %d got %d", want, order.TotalCents) }
}
