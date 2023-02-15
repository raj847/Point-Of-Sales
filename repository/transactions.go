package repository

import (
	"encoding/json"
	"vandesar/entity"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	AddTrans(trans entity.TransactionReq) []error
	UpdateTrans(trans entity.Transaction) error
	DeleteTrans(id int) error
	ReadTrans() ([]entity.TransactionReq, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

// func (c *transactionRepository) AddTrans(prods []entity.Prods, money float64, quantity float64, totalharga float64, catatan string) error {
func (c *transactionRepository) AddTrans(trans entity.TransactionReq) []error {
	errs := []error{}

	err := c.db.Transaction(func(tx *gorm.DB) error {
		products, _ := json.Marshal(trans.Products)
		transaction := entity.Transaction{
			UserID:     trans.UserID,
			Debt:       trans.Debt,
			Status:     trans.Status,
			Money:      trans.Money,
			Products:   products,
			TotalHarga: trans.TotalHarga,
			Notes:      trans.Notes,
			TotalLaba:  trans.TotalLaba,
		}

		err := tx.Create(&transaction).Error
		if err != nil {
			return err
		}

		for _, v := range trans.Products {
			// get product
			var product entity.Product
			err = tx.Table("products").Where("id = ?", v.ProductID).First(&product).Error
			if err != nil {
				return err
			}

			// validate stock
			//if product.Stock < v.Quantity {
			//	errs = append(errs, fmt.Errorf("stock for product id %d not enough", product.ID))
			//	continue
			//}

			// update
			id := v.ProductID
			err = tx.Table("products").Where("id = ?", id).Update("stock", gorm.Expr("stock - ?", v.Quantity)).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		errs = append(errs, err)
	}

	return errs // TODO: replace this
}

func (c *transactionRepository) UpdateTrans(trans entity.Transaction) error {
	c.db.Table("transactions").Where("id = ?", trans.ID).Updates(&trans)
	return nil // TODO: replace this
}

func (c *transactionRepository) DeleteTrans(id int) error {
	err := c.db.Delete(&entity.Transaction{}, id).Error
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (c *transactionRepository) ReadTrans() ([]entity.TransactionReq, error) {
	transactions := []entity.Transaction{}
	c.db.Table("transactions").Select("*").Scan(&transactions)

	resp := make([]entity.TransactionReq, 0, len(transactions))

	for _, v := range transactions {
		var products []entity.Prods
		json.Unmarshal(v.Products, &products)

		resp = append(resp, entity.TransactionReq{
			UserID:     v.UserID,
			Debt:       v.Debt,
			Status:     v.Status,
			Money:      v.Money,
			Products:   products,
			TotalHarga: v.TotalHarga,
			Notes:      v.Notes,
			TotalLaba:  v.TotalLaba,
		})
	}

	return resp, nil
}
