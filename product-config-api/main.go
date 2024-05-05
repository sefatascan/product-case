package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
	"product-config-api/internal/client"
	"product-config-api/internal/config"
	"product-config-api/internal/controller"
	"product-config-api/internal/repository"
	"product-config-api/internal/service"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	_ = container.Provide(config.NewApplicationConfigManager)
	_ = container.Provide(client.NewRedisClient)
	_ = container.Provide(repository.NewProductConfigRepository)
	_ = container.Provide(service.NewProductConfigService)
	_ = container.Provide(controller.NewProductConfigController)

	return container
}

func BuildApp(container *dig.Container) *fiber.App {
	app := fiber.New()

	err := container.Invoke(func(productController controller.ProductConfigController) {
		app.Get("/product/config", productController.GetProductConfig)
		app.Post("/product/config", productController.CreateProductConfig)
	})

	if err != nil {
		panic(err)
	}

	return app
}

func main() {
	container := BuildContainer()
	app := BuildApp(container)

	err := container.Invoke(func(applicationConfig config.ApplicationConfigManager) {
		_ = app.Listen(fmt.Sprintf(":%d", applicationConfig.Server.Port))
	})

	if err != nil {
		panic(err)
	}
}
