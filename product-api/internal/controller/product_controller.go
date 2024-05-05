package controller

import (
	"github.com/gofiber/fiber/v2"
	"product-api/internal/model"
	"product-api/internal/service"
)

type ProductController interface {
	GetProducts(c *fiber.Ctx) error
}

type ProductControllerImpl struct {
	productService service.ProductService
}

func (p *ProductControllerImpl) GetProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 0)
	if page < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'page' must be greater than 0",
		})
	}

	size := c.QueryInt("size", 10)
	if size <= 0 || size > 10000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'size' must be between 0 and 10000",
		})
	}

	filterModel := new(model.ProductRequestFilterModel)
	if err := c.QueryParser(filterModel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'filters' invalid",
		})
	}

	products, err := p.productService.GetProducts(page, size, *filterModel)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"page": page,
		"size": size,
		"data": fiber.Map{
			"products": products,
		},
	})
}

func NewProductController(productService service.ProductService) ProductController {
	return &ProductControllerImpl{
		productService: productService,
	}
}
