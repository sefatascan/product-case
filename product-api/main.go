package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
	"product-api/internal/cache"
	"product-api/internal/client"
	"product-api/internal/config"
	"product-api/internal/controller"
	"product-api/internal/service"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	_ = container.Provide(config.NewApplicationConfigManager)
	_ = container.Provide(client.NewProductElasticClient)
	_ = container.Provide(client.NewProductConfigClient)
	_ = container.Provide(cache.NewProductConfigCache)
	_ = container.Provide(client.NewRedisClient)
	_ = container.Provide(service.NewFilterService)
	_ = container.Provide(service.NewElasticsearchQueryService)
	_ = container.Provide(service.NewProductService)
	_ = container.Provide(controller.NewProductController)

	return container
}

func BuildApp(container *dig.Container) *fiber.App {
	app := fiber.New()

	err := container.Invoke(func(productController controller.ProductController) {
		app.Get("/products", productController.GetProducts)
	})

	if err != nil {
		panic(err)
	}

	return app
}

func BuildSubscriber(container *dig.Container) {
	err := container.Invoke(func(productConfigCache cache.ProductConfigCache, redisClient client.RedisClient) {
		redisClient.Subscribe(func(payload string) {
			if err := productConfigCache.ReloadCache(); err != nil {
				fmt.Printf("cache invalidate event not recived of time:%s", payload)
			}
		})
	})

	if err != nil {
		panic(err)
	}
}

func main() {

	container := BuildContainer()
	app := BuildApp(container)
	BuildSubscriber(container)

	err := container.Invoke(func(applicationConfig config.ApplicationConfigManager) {
		_ = app.Listen(fmt.Sprintf(":%d", applicationConfig.Server.Port))
	})

	if err != nil {
		panic(err)
	}
}
