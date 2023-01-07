package service

import (
	"context"
	"vandesar/entity"
	"vandesar/repository"
)

type ProductService interface {
	GetProducts(ctx context.Context, id int) ([]entity.Product, error)
	AddProduct(ctx context.Context, product *entity.Product) (entity.Product, error)
	GetProductByID(ctx context.Context, id int) (entity.Product, error)
	UpdateProduct(ctx context.Context, Product *entity.Product) (entity.Product, error)
	DeleteProduct(ctx context.Context, id int) error
	GetProductBySearch(ctx context.Context, name string) ([]entity.Product, error)
}

type productService struct {
	prodRepo repository.ProductRepository
}

func NewProductService(prodRepo repository.ProductRepository) ProductService {
	return &productService{prodRepo}
}

func (s *productService) GetProducts(ctx context.Context, id int) ([]entity.Product, error) {
	return s.prodRepo.GetProductsByUserId(ctx, id)
}

func (s *productService) AddProduct(ctx context.Context, product *entity.Product) (entity.Product, error) {
	err := s.prodRepo.AddProduct(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}
	return *product, nil
}

func (s *productService) GetProductByID(ctx context.Context, id int) (entity.Product, error) {
	return s.prodRepo.GetProductByID(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, product *entity.Product) (entity.Product, error) {
	err := s.prodRepo.UpdateProduct(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}
	return *product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	return s.prodRepo.DeleteProduct(ctx, id)
}

func (s *productService) GetProductBySearch(ctx context.Context, name string) ([]entity.Product, error) {
	return s.prodRepo.GetProductBySearch(ctx, name)
}
