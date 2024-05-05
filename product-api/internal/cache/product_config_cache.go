package cache

import (
	"fmt"
	"product-api/internal/client"
	"product-api/internal/model"
)

type ProductConfigCache interface {
	GetCache() model.ProductConfigModel
	ReloadCache() error
}

type ProductConfigCacheImpl struct {
	productConfigClient client.ProductConfigClient
	productConfig       model.ProductConfigModel
}

func (p *ProductConfigCacheImpl) GetCache() model.ProductConfigModel {
	return p.productConfig
}

func (p *ProductConfigCacheImpl) ReloadCache() error {
	productConfigModel, err := p.productConfigClient.GetConfig()
	if err != nil {
		return fmt.Errorf("product config cache not reloaded : %s", err)
	}
	p.productConfig = productConfigModel
	return err
}

func NewProductConfigCache(productConfigClient client.ProductConfigClient) ProductConfigCache {
	productConfigCache := &ProductConfigCacheImpl{
		productConfigClient: productConfigClient,
	}
	_ = productConfigCache.ReloadCache()
	return productConfigCache
}
