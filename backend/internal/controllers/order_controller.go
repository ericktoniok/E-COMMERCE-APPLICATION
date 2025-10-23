package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"mini-ecommerce/backend/internal/appctx"
	"mini-ecommerce/backend/internal/middlewares"
	"mini-ecommerce/backend/internal/services"
)

type OrderController struct {
	Ax  *appctx.Context
	Svc *services.OrderService
}

type checkoutReq struct {
	Items []services.CheckoutItem `json:"items"`
}

type checkoutResp struct {
	OrderID       uint   `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

func NewOrderController(ax *appctx.Context, svc *services.OrderService) *OrderController {
	return &OrderController{Ax: ax, Svc: svc}
}

func (h *OrderController) Checkout(c *fiber.Ctx) error {
	var req checkoutReq
	if err := c.BodyParser(&req); err != nil || len(req.Items) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	uid, _ := c.Locals(middlewares.CtxUserID).(uint)
	order, tx, err := h.Svc.Checkout(uid, req.Items)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(checkoutResp{OrderID: order.ID, TransactionID: tx.ProviderRef, Status: string(tx.Status)})
}

func (h *OrderController) ListMine(c *fiber.Ctx) error {
	uid, _ := c.Locals(middlewares.CtxUserID).(uint)
	orders, err := h.Svc.ListByUser(uid)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

func (h *OrderController) AdminList(c *fiber.Ctx) error {
	orders, err := h.Svc.ListAll()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}
