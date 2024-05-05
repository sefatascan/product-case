package service

import (
	"fmt"
	"product-api/internal/model"
)

type ProductService interface {
	GetProducts(page int, size int, filters model.ProductRequestFilterModel) (interface{}, error)
}

type ProductServiceImpl struct {
	filterService       FilterService
	elasticQueryService ElasticsearchQueryService
}

func (p *ProductServiceImpl) GetProducts(page int, size int, filters model.ProductRequestFilterModel) (interface{}, error) {
	elasticFilters, err := p.filterService.ValidateAndConvertFilters(filters)
	if err != nil {
		return nil, fmt.Errorf("not supported filters")
	}

	elasticQuery, err := p.elasticQueryService.CreateElasticQuery(page, size, elasticFilters)
	if err != nil {
		return nil, fmt.Errorf("elastic query creation error")
	}

	products, err := p.elasticQueryService.GetProducts(elasticQuery)
	if err != nil {
		return nil, fmt.Errorf("elastic client error %v", err)
	}

	return products, nil
}

func NewProductService(filterService FilterService, elasticQueryService ElasticsearchQueryService) ProductService {
	return &ProductServiceImpl{
		filterService:       filterService,
		elasticQueryService: elasticQueryService,
	}
}
