package repository

import (
	"context"
	"fmt"
	"vandesar/entity"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductsByUserId(ctx context.Context, userId, adminId int) ([]entity.Product, error)
	AddProduct(ctx context.Context, products *entity.Product) error
	GetProductByID(ctx context.Context, id int) (entity.Product, error)
	DeleteProduct(ctx context.Context, id int) error
	UpdateProduct(ctx context.Context, products *entity.Product) error
	GetProductBySearch(ctx context.Context, object string) ([]entity.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (p *productRepository) GetProductsByUserId(ctx context.Context, userId, adminId int) ([]entity.Product, error) {
	res := []entity.Product{}
	cek, err := p.db.WithContext(ctx).Table("products").Select("*").Where("admin_id = ? ", adminId).Rows()
	if err != nil {
		return []entity.Product{}, err
	}
	defer cek.Close()
	for cek.Next() {
		p.db.ScanRows(cek, &res)
	}

	return res, nil
	// TODO: replace this
}

func (p *productRepository) AddProduct(ctx context.Context, products *entity.Product) error {
	err := p.db.WithContext(ctx).Create(&products).Error
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (p *productRepository) GetProductByID(ctx context.Context, id int) (entity.Product, error) {
	res := entity.Product{}
	err := p.db.WithContext(ctx).Table("products").Where("id = ?", id).Find(&res).Error

	if err != nil {
		return entity.Product{}, err
	}

	return res, nil // TODO: replace this
}

func (p *productRepository) DeleteProduct(ctx context.Context, id int) error {
	err := p.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (p *productRepository) UpdateProduct(ctx context.Context, products *entity.Product) error {
	err := p.db.WithContext(ctx).Table("products").Where("id = ?", products.ID).Updates(&products).Error
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (p *productRepository) GetProductBySearch(ctx context.Context, object string) ([]entity.Product, error) {
	res := []entity.Product{}
	err := p.db.WithContext(ctx).Table("products").Where("name LIKE ?", fmt.Sprintf("%s%s%s", "%", object, "%")).Find(&res).Error

	if err != nil {
		return []entity.Product{}, err
	}

	return res, nil // TODO: replace this
}
