package service

import (
	"context"
	"errors"
	"vandesar/entity"
	"vandesar/repository"
)

type ProductService struct {
	prodRepo *repository.ProductRepository
}

func NewProductService(prodRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		prodRepo: prodRepo,
	}
}

func (s *ProductService) GetProducts(ctx context.Context, adminId uint) ([]entity.Product, error) {
	return s.prodRepo.GetProductsByUserId(ctx, adminId)
}

func (s *ProductService) AddProduct(ctx context.Context, product entity.Product) (entity.Product, error) {
	err := s.prodRepo.AddProduct(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (entity.Product, error) {
	return s.prodRepo.GetProductByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product entity.Product) (entity.Product, error) {
	existingProduct, err := s.GetProductByID(ctx, int(product.ID))
	if err != nil {
		return entity.Product{}, err
	}

	if existingProduct.UserID != product.UserID {
		return entity.Product{}, errors.New("you are not allowed to update this product")
	}

	err = s.prodRepo.UpdateProduct(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return s.prodRepo.DeleteProduct(ctx, id)
}

func (s *ProductService) SearchProducts(ctx context.Context, name string) ([]entity.Product, error) {
	return s.prodRepo.SearchProducts(ctx, name)
}
