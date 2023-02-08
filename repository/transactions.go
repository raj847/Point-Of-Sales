package repository

import (
	"vandesar/entity"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	AddTrans(trans entity.Transaction) error
	UpdateTrans(id int, trans entity.Transaction) error
	DeleteTrans(id int) error
	ReadTrans() ([]entity.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

// func (c *transactionRepository) AddTrans(prods []entity.Prods, money float64, quantity float64, totalharga float64, catatan string) error {
func (c *transactionRepository) AddTrans(trans entity.Transaction) error {

	err := c.db.Transaction(func(tx *gorm.DB) error {
		product := entity.Product{}
		for _, v := range trans.Products {
			trans.TotalHarga += v.TotalPrice
		}
		trans.Debt = trans.TotalHarga - trans.Money
		if trans.Debt != 0 {
			trans.Status = "Lunas"
		} else {
			trans.Status = "Hutang"
		}
		err := tx.Create(&trans).Error
		if err != nil {
			return err
		}

		for _, v := range trans.Products {
			if v.ProductID == product.ID {
				err = tx.Table("products").Update("stock", product.Stock-v.Quantity).Error
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err // TODO: replace this
}

func (c *transactionRepository) UpdateTrans(id int, trans entity.Transaction) error {
	c.db.Table("transactions").Where("id = ?", id).Updates(&trans)
	return nil // TODO: replace this
}

func (c *transactionRepository) DeleteTrans(id int) error {
	err := c.db.Delete(&entity.Transaction{}, id).Error
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (c *transactionRepository) ReadTrans() ([]entity.Transaction, error) {
	BacaTransaksi := []entity.Transaction{}
	c.db.Table("transactions").Select("transactions.id, transactions.product_id,products.name, carts.quantity,carts.total_price ").Joins("inner join bla bla bla where carts.deleted_at IS NULL").Scan(&BacaTransaksi)
	return BacaTransaksi, nil
}
