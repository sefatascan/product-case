package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"product-api/internal/config"
	"product-api/internal/model"
)

type ProductConfigClient interface {
	GetConfig() (model.ProductConfigModel, error)
}

type ProductConfigClientImpl struct {
	applicationConfig config.ApplicationConfigManager
}

func (p *ProductConfigClientImpl) GetConfig() (model.ProductConfigModel, error) {
	response, err := http.Get(fmt.Sprintf("%s/product/config", p.applicationConfig.ProductConfigClient.Url))
	if err != nil {
		return model.ProductConfigModel{}, fmt.Errorf("http error when fetch config err : %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.ProductConfigModel{}, fmt.Errorf("io error when read config err : %s", err)
	}

	var productConfigModel model.ProductConfigModel
	err = json.Unmarshal(body, &productConfigModel)
	if err != nil {
		return model.ProductConfigModel{}, fmt.Errorf("unmarshal error when deserialize config err : %s", err)
	}

	return productConfigModel, nil

}

func NewProductConfigClient(applicationConfig config.ApplicationConfigManager) ProductConfigClient {
	return &ProductConfigClientImpl{
		applicationConfig: applicationConfig,
	}
}
