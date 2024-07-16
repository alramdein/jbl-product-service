package usecase

import (
	"context"
	"errors"
	"product-service/model"
	"product-service/repository"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

type productUsecase struct {
	productRepo repository.IProductRepository
}

func NewProductUsecase(repo repository.IProductRepository) IProductUsecase {
	return &productUsecase{repo}
}

func (u *productUsecase) GetProducts(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	products, err := u.productRepo.GetProducts(ctx, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return products, nil
}

func (u *productUsecase) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}
	return u.productRepo.GetProductByID(ctx, id)
}

func (u *productUsecase) CreateProduct(ctx context.Context, product model.Product) error {
	existingProduct, err := u.productRepo.GetProductBySKU(ctx, product.SKU)
	if err == nil && existingProduct.ID != "" {
		return errors.New("product with this SKU already exists")
	}

	product.ID = uuid.NewString()
	return u.productRepo.CreateProduct(ctx, product)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, product model.Product) error {
	_, err := uuid.Parse(product.ID)
	if err != nil {
		return errors.New("invalid product ID")
	}
	return u.productRepo.UpdateProduct(ctx, product)
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid product ID")
	}
	return u.productRepo.DeleteProduct(ctx, id)
}
