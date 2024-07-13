package repository

import (
	"context"
	"database/sql"
	"product-service/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) GetProducts(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	rows, err := r.db.Query("SELECT id, name, sku, image, price FROM products LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.SKU, &product.Image, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	err := r.db.QueryRow("SELECT id, name, sku, image, price, description FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.SKU, &product.Image, &product.Price, &product.Description)
	if err != nil {
		return &product, err
	}
	return &product, nil
}

func (r *ProductRepository) GetProductBySKU(ctx context.Context, sku string) (*model.Product, error) {
	var product model.Product
	err := r.db.QueryRow("SELECT id, name, sku, image, price, description FROM products WHERE sku = $1", sku).
		Scan(&product.ID, &product.Name, &product.SKU, &product.Image, &product.Price, &product.Description)
	if err != nil {
		return &product, err
	}
	return &product, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product model.Product) error {
	_, err := r.db.Exec("INSERT INTO products (name, sku, image, price, description) VALUES ($1, $2, $3, $4, $5)",
		product.Name, product.SKU, product.Image, product.Price, product.Description)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product model.Product) error {
	_, err := r.db.Exec("UPDATE products SET name = $1, sku = $2, image = $3, price = $4, description = $5 WHERE id = $6",
		product.Name, product.SKU, product.Image, product.Price, product.Description, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
