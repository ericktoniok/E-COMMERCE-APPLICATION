package controllers

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"mini-ecommerce/backend/internal/appctx"
	"mini-ecommerce/backend/internal/models"
	"mini-ecommerce/backend/internal/services"
)

type ProductController struct {
	Ax   *appctx.Context
	Svc  *services.ProductService
}

func NewProductController(ax *appctx.Context, svc *services.ProductService) *ProductController {
	return &ProductController{Ax: ax, Svc: svc}
}

func (h *ProductController) List(c *fiber.Ctx) error {
	items, err := h.Svc.List()
	if err != nil { return c.Status(500).JSON(fiber.Map{"error":"failed"}) }
	return c.JSON(items)
}

func (h *ProductController) Create(c *fiber.Ctx) error {
	var p models.Product
	if err := c.BodyParser(&p); err != nil { return c.Status(400).JSON(fiber.Map{"error":"bad request"}) }
	if strings.TrimSpace(p.Name) == "" || p.PriceCents <= 0 { return c.Status(400).JSON(fiber.Map{"error":"invalid payload"}) }
	if err := h.Svc.Create(&p); err != nil { return c.Status(400).JSON(fiber.Map{"error": err.Error()}) }
	return c.Status(201).JSON(p)
}

func (h *ProductController) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var p models.Product
	if err := c.BodyParser(&p); err != nil { return c.Status(400).JSON(fiber.Map{"error":"bad request"}) }
	p.ID = uint(id)
	if err := h.Svc.Update(&p); err != nil { return c.Status(400).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(p)
}

func (h *ProductController) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.Svc.Delete(uint(id)); err != nil { return c.Status(400).JSON(fiber.Map{"error": err.Error()}) }
	return c.SendStatus(204)
}

func (h *ProductController) UploadImage(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := c.FormFile("image")
	if err != nil { return c.Status(400).JSON(fiber.Map{"error":"image required"}) }
	ext := filepath.Ext(file.Filename)
	name := fmt.Sprintf("p_%s_%d%s", id, time.Now().Unix(), ext)
	path := filepath.Join(h.Ax.Cfg.ImageStoragePath, name)
	if err := c.SaveFile(file, path); err != nil { return c.Status(500).JSON(fiber.Map{"error":"save failed"}) }
	imageURL := "/images/" + name
	// update product
	pid, _ := strconv.Atoi(id)
	p := models.Product{ID: uint(pid), ImageURL: imageURL}
	if err := h.Svc.Update(&p); err != nil { return c.Status(500).JSON(fiber.Map{"error":"update failed"}) }
	return c.JSON(fiber.Map{"image_url": imageURL})
}
