package service

import (
	"product-api/internal/cache"
	"product-api/internal/client"
)

type ElasticsearchQueryService interface {
	CreateElasticQuery(page int, size int, elasticFilters map[string]filters) (map[string]interface{}, error)
	GetProducts(query map[string]interface{}) (interface{}, error)
}

type ElasticsearchQueryServiceImpl struct {
	productElasticClient      client.ProductElasticClient
	productElasticConfigCache cache.ProductConfigCache
}

func (e *ElasticsearchQueryServiceImpl) GetProducts(query map[string]interface{}) (interface{}, error) {
	return e.productElasticClient.GetProducts(query)
}

func (e *ElasticsearchQueryServiceImpl) CreateElasticQuery(page int, size int, elasticFilters map[string]filters) (map[string]interface{}, error) {
	productConfig := e.productElasticConfigCache.GetCache()

	var rangeQueries []map[string]interface{}

	for field, filter := range elasticFilters {
		rangeQuery := map[string]interface{}{
			"range": map[string]interface{}{
				field: filter,
			},
		}
		rangeQueries = append(rangeQueries, rangeQuery)
	}

	query := map[string]interface{}{
		"from": page * size,
		"size": size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": rangeQueries,
			},
		},
		"sort": map[string]interface{}{
			productConfig.FieldName: map[string]interface{}{
				"order": productConfig.OrderType,
			},
		},
	}

	return query, nil
}

func NewElasticsearchQueryService(productElasticClient client.ProductElasticClient, productElasticConfigCache cache.ProductConfigCache) ElasticsearchQueryService {
	return &ElasticsearchQueryServiceImpl{
		productElasticClient:      productElasticClient,
		productElasticConfigCache: productElasticConfigCache,
	}
}
