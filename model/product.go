package model

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}
