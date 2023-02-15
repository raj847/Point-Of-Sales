package service

import (
	"context"
	"vandesar/entity"
)

type productRepository interface {
	GetProductsByUserId(ctx context.Context, adminId uint) ([]entity.Product, error)
	AddProduct(ctx context.Context, products *entity.Product) error
	GetProductByID(ctx context.Context, id int) (entity.Product, error)
	DeleteProduct(ctx context.Context, id int) error
	UpdateProduct(ctx context.Context, products *entity.Product) error
	GetProductBySearch(ctx context.Context, object string) ([]entity.Product, error)
}

type ProductService struct {
	prodRepo productRepository
}

func NewProductService(prodRepo productRepository) *ProductService {
	return &ProductService{
		prodRepo: prodRepo,
	}
}

func (s *ProductService) GetProducts(ctx context.Context, adminId uint) ([]entity.Product, error) {
	return s.prodRepo.GetProductsByUserId(ctx, adminId)
}

func (s *ProductService) AddProduct(ctx context.Context, product *entity.Product) (entity.Product, error) {
	err := s.prodRepo.AddProduct(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}
	return *product, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (entity.Product, error) {
	return s.prodRepo.GetProductByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *entity.Product) (entity.Product, error) {
	err := s.prodRepo.UpdateProduct(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}
	return *product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return s.prodRepo.DeleteProduct(ctx, id)
}

func (s *ProductService) GetProductBySearch(ctx context.Context, name string) ([]entity.Product, error) {
	return s.prodRepo.GetProductBySearch(ctx, name)
}
