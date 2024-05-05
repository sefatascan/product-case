package model

import "gorm.io/gorm"

type ProductConfigEntity struct {
	gorm.Model
	FieldName string `json:"field_name"`
	OrderType string `json:"order_type"`
	Active    bool   `json:"active"`
}
