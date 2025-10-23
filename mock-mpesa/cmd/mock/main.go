package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type payReq struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount_cents"`
}

type payResp struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

func main() {
	port := os.Getenv("MOCK_PORT")
	if port == "" { port = "8090" }
	webhook := os.Getenv("BACKEND_WEBHOOK_URL")
	if webhook == "" { webhook = "http://api:8080/api/webhooks/mpesa" }

	app := fiber.New()

	app.Post("/pay", func(c *fiber.Ctx) error {
		var req payReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "bad request"})
		}
		txID := time.Now().Format("20060102150405")
		go func(orderID string, amount int64, tx string) {
			time.Sleep(3 * time.Second)
			payload := map[string]any{
				"transaction_id": tx,
				"order_id": orderID,
				"status": "SUCCESS",
				"amount_cents": amount,
			}
			b, _ := json.Marshal(payload)
			resp, err := http.Post(webhook, "application/json", bytesReader(b))
			if err != nil {
				log.Println("webhook error:", err)
				return
			}
			resp.Body.Close()
		}(req.OrderID, req.Amount, txID)

		return c.JSON(payResp{TransactionID: txID, Status: "PENDING"})
	})

	log.Printf("Mock M-Pesa listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

// bytesReader is a tiny helper to avoid importing bytes directly in multiple places
func bytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }
