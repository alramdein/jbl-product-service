package repository

import (
	"context"
	"product-service/model"
)

type IProductRepository interface {
	GetProducts(ctx context.Context, limit, offset int) ([]*model.Product, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*model.Product, error)
	CreateProduct(ctx context.Context, product model.Product) error
	UpdateProduct(ctx context.Context, product model.Product) error
	DeleteProduct(ctx context.Context, id string) error
}
