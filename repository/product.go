package repository

import (
	"context"
	"fmt"
	"vandesar/entity"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (p *ProductRepository) GetProductsByUserId(ctx context.Context, userId uint) ([]entity.Product, error) {
	var productsResult []entity.Product

	products, err := p.db.
		WithContext(ctx).
		Table("products").
		Select("*").
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Rows()
	if err != nil {
		return []entity.Product{}, err
	}
	defer products.Close()

	for products.Next() {
		p.db.ScanRows(products, &productsResult)
	}

	return productsResult, nil
}

func (p *ProductRepository) AddProduct(ctx context.Context, products entity.Product) error {
	err := p.db.
		WithContext(ctx).
		Create(&products).Error
	return err
}

func (p *ProductRepository) AddProductAkeh(ctx context.Context, products []entity.Product) error {
	err := p.db.
		WithContext(ctx).
		Create(&products).Error
	return err
}

func (p *ProductRepository) GetProductByID(ctx context.Context, id int) (entity.Product, error) {
	var productResult entity.Product

	err := p.db.
		WithContext(ctx).
		Table("products").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&productResult).Error
	if err != nil {
		return entity.Product{}, err
	}

	return productResult, nil
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	err := p.db.
		WithContext(ctx).
		Delete(&entity.Product{}, id).Error
	return err
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, products entity.Product) error {
	err := p.db.
		WithContext(ctx).
		Table("products").
		Where("id = ?", products.ID).
		Updates(&products).Error
	return err
}

func (p *ProductRepository) SearchProducts(ctx context.Context, object string) ([]entity.Product, error) {
	var productResult []entity.Product

	err := p.db.
		WithContext(ctx).
		Table("products").
		Where("name LIKE ? AND deleted_at IS NULL", fmt.Sprintf("%s%s%s", "%", object, "%")).
		Find(&productResult).Error
	if err != nil {
		return []entity.Product{}, err
	}
	return productResult, nil
}
