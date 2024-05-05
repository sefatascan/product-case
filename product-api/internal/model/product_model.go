package model

type ProductModel struct {
	Name     string  `json:"name"`
	ItemId   string  `json:"item_id"`
	Locale   string  `json:"locale"`
	Click    float64 `json:"click"`
	Purchase float64 `json:"purchase"`
}
