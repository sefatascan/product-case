package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"product-api/internal/config"
	"product-api/internal/model"
)

type ProductElasticClient interface {
	GetProducts(query map[string]interface{}) (interface{}, error)
}

type ProductElasticClientImpl struct {
	applicationConfig config.ApplicationConfigManager
	elasticClient     *elasticsearch.Client
}

func (p *ProductElasticClientImpl) GetProducts(query map[string]interface{}) (interface{}, error) {

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error elastic client %v", err)
	}

	res, err := p.elasticClient.Search(
		p.elasticClient.Search.WithContext(context.Background()),
		p.elasticClient.Search.WithIndex(p.applicationConfig.Elasticsearch.Index),
		p.elasticClient.Search.WithBody(&buf),
		p.elasticClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error performing search request: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		}
		return nil, fmt.Errorf("elasticsearch error: %s", e)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}
	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("hits not found in Elasticsearch response")
	}
	var products []model.ProductModel
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})

		product := model.ProductModel{
			Name:     source["name"].(string),
			ItemId:   source["item_id"].(string),
			Locale:   source["locale"].(string),
			Click:    source["click"].(float64),
			Purchase: source["purchase"].(float64),
		}

		products = append(products, product)
	}

	return products, nil

}

func NewProductElasticClient(applicationConfig config.ApplicationConfigManager) ProductElasticClient {
	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("%s:%d", applicationConfig.Elasticsearch.Host, applicationConfig.Elasticsearch.Port),
		},
	}
	elasticClient, err := elasticsearch.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	return &ProductElasticClientImpl{
		applicationConfig: applicationConfig,
		elasticClient:     elasticClient,
	}
}
