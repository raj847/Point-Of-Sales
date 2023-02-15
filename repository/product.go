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

func (p *ProductRepository) GetProductsByUserId(ctx context.Context, adminId uint) ([]entity.Product, error) {
	var res []entity.Product
	cek, err := p.db.WithContext(ctx).Table("products").Select("*").Where("user_id = ? ", adminId).Rows()
	if err != nil {
		return []entity.Product{}, err
	}
	defer cek.Close()
	for cek.Next() {
		p.db.ScanRows(cek, &res)
	}

	return res, nil
}

func (p *ProductRepository) AddProduct(ctx context.Context, products *entity.Product) error {
	err := p.db.WithContext(ctx).Create(&products).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) GetProductByID(ctx context.Context, id int) (entity.Product, error) {
	res := entity.Product{}
	err := p.db.WithContext(ctx).Table("products").Where("id = ?", id).Find(&res).Error

	if err != nil {
		return entity.Product{}, err
	}

	return res, nil
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	err := p.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, products *entity.Product) error {
	err := p.db.WithContext(ctx).Table("products").Where("id = ?", products.ID).Updates(&products).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) GetProductBySearch(ctx context.Context, object string) ([]entity.Product, error) {
	var res []entity.Product
	err := p.db.WithContext(ctx).Table("products").Where("name LIKE ?", fmt.Sprintf("%s%s%s", "%", object, "%")).Find(&res).Error

	if err != nil {
		return []entity.Product{}, err
	}

	return res, nil
}
