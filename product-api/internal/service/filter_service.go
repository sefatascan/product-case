package service

import (
	"fmt"
	"product-api/internal/model"
	"regexp"
	"strconv"
)

var (
	regexPattern = `^(gte|lte|lt|gt|eq):(\d+)$`
	regex        = regexp.MustCompile(regexPattern)
)

type filters map[string]int

type FilterService interface {
	ValidateAndConvertFilters(requestFilters model.ProductRequestFilterModel) (map[string]filters, error)
}

type FilterServiceImpl struct {
}

func (f FilterServiceImpl) ValidateAndConvertFilters(requestFilters model.ProductRequestFilterModel) (map[string]filters, error) {

	elasticFilters := make(map[string]filters)

	itemIdFilters, err := f.compileRegex(requestFilters.ItemId)
	if err != nil {
		return nil, err
	}
	elasticFilters["item_id"] = itemIdFilters

	clickCountFilters, err := f.compileRegex(requestFilters.ClickCount)
	if err != nil {
		return nil, err
	}
	elasticFilters["click"] = clickCountFilters

	purchaseCountFilters, err := f.compileRegex(requestFilters.PurchaseCount)
	if err != nil {
		return nil, err
	}
	elasticFilters["purchase"] = purchaseCountFilters

	return elasticFilters, nil
}

func (f FilterServiceImpl) compileRegex(values []string) (filters, error) {
	filterMap := make(map[string]int)

	for _, value := range values {
		matches := regex.FindStringSubmatch(value)

		if len(matches) != 3 {
			return nil, fmt.Errorf("unsupported filter pattern")
		}

		filter := matches[1]
		number, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("unsupported filter value %s:%d", filter, number)

		}
		filterMap[filter] = number
	}
	return filterMap, nil
}

func NewFilterService() FilterService {
	return &FilterServiceImpl{}
}
