package model

type ProductRequestFilterModel struct {
	ItemId        []string `query:"item_id"`
	ClickCount    []string `query:"click"`
	PurchaseCount []string `query:"purchase"`
}
