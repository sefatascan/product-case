package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"product-config-api/internal/model/request"
	"product-config-api/internal/service"
)

type ProductConfigController interface {
	CreateProductConfig(c *fiber.Ctx) error
	GetProductConfig(c *fiber.Ctx) error
}

type ProductConfigControllerImpl struct {
	productConfigService service.ProductConfigService
}

func (p *ProductConfigControllerImpl) CreateProductConfig(c *fiber.Ctx) error {

	requestBody := c.Body()
	var productConfig request.ProductConfigRequestModel

	if err := json.Unmarshal(requestBody, &productConfig); err != nil {
		return err
	}

	err := p.productConfigService.CreateProductConfig(productConfig)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("product config not created, %v", err),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (p *ProductConfigControllerImpl) GetProductConfig(c *fiber.Ctx) error {
	productConfig, err := p.productConfigService.GetActiveProductConfig()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("product config not found, %v", err),
		})
	}
	return c.Status(fiber.StatusOK).JSON(productConfig)
}

func NewProductConfigController(productConfigService service.ProductConfigService) ProductConfigController {
	return &ProductConfigControllerImpl{
		productConfigService: productConfigService,
	}
}
