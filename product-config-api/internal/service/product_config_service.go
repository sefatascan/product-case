package service

import (
	"product-config-api/internal/client"
	"product-config-api/internal/model"
	"product-config-api/internal/model/request"
	"product-config-api/internal/repository"
)

type ProductConfigService interface {
	CreateProductConfig(model request.ProductConfigRequestModel) error
	GetActiveProductConfig() (model.ProductConfigEntity, error)
}

type ProductConfigServiceImpl struct {
	productConfigRepository repository.ProductConfigRepository
	redisClient             client.RedisClient
}

func (p *ProductConfigServiceImpl) CreateProductConfig(model request.ProductConfigRequestModel) error {
	err := p.productConfigRepository.CreateProductConfig(model)
	if err != nil {
		return err
	}
	_ = p.redisClient.PublishEvent()
	return nil
}

func (p *ProductConfigServiceImpl) GetActiveProductConfig() (model.ProductConfigEntity, error) {
	return p.productConfigRepository.GetActiveProductConfig()
}

func NewProductConfigService(productConfigRepository repository.ProductConfigRepository, redisClient client.RedisClient) ProductConfigService {
	return &ProductConfigServiceImpl{
		productConfigRepository: productConfigRepository,
		redisClient:             redisClient,
	}
}
