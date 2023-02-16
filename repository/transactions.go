package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"vandesar/entity"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (c *transactionRepository) AddTrans(ctx context.Context, trans entity.TransactionReq) []error {
	var errs []error

	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cartList, _ := json.Marshal(trans.CartList)

		transaction := entity.Transaction{
			UserID:      trans.UserID,
			Debt:        trans.Debt,
			Status:      trans.Status,
			Money:       trans.Money,
			CartList:    cartList,
			TotalPrice:  trans.TotalPrice,
			Notes:       trans.Notes,
			TotalProfit: trans.TotalProfit,
		}

		err := tx.Create(&transaction).Error
		if err != nil {
			return err
		}

		for _, v := range trans.CartList {
			var product entity.Product
			err = tx.Table("products").Where("id = ?", v.ProductID).First(&product).Error
			if err != nil {
				return err
			}

			if product.Stock < v.Quantity {
				errs = append(errs, fmt.Errorf("stock for product id %d not enough", product.ID))
				continue
			}

			err = tx.Table("products").Where("id = ?", v.ProductID).Update("stock", gorm.Expr("stock - ?", v.Quantity)).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (c *transactionRepository) UpdateTrans(trans entity.Transaction) error {
	c.db.Table("transactions").Where("id = ?", trans.ID).Updates(&trans)
	return nil
}

func (c *transactionRepository) DeleteTrans(id uint) error {
	err := c.db.Delete(&entity.Transaction{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *transactionRepository) ReadTransByCashier(userId uint) ([]entity.TransactionReq, error) {
	var transactions []entity.Transaction
	err := c.db.
		Table("transactions").
		Select("*").
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Scan(&transactions).Error

	if err != nil {
		return nil, err
	}

	resp := make([]entity.TransactionReq, 0, len(transactions))

	for _, v := range transactions {
		var cartList []entity.Prods
		_ = json.Unmarshal(v.CartList, &cartList)

		resp = append(resp, entity.TransactionReq{
			UserID:      v.UserID,
			Debt:        v.Debt,
			Status:      v.Status,
			Money:       v.Money,
			CartList:    cartList,
			TotalPrice:  v.TotalPrice,
			Notes:       v.Notes,
			TotalProfit: v.TotalProfit,
		})
	}

	return resp, nil
}

func (c *transactionRepository) ReadTransByAdmin(adminId int) ([]entity.TransactionReq, error) {
	var transactions []entity.Transaction
	err := c.db.
		Table("transactions").
		Select("*").
		Joins("JOIN users ON users.id = transactions.user_id").
		Where("users.admin_id = ?", adminId).
		Scan(&transactions).Error
	if err != nil {
		return nil, err
	}

	resp := make([]entity.TransactionReq, 0, len(transactions))

	for _, v := range transactions {
		var cartList []entity.Prods
		_ = json.Unmarshal(v.CartList, &cartList)

		resp = append(resp, entity.TransactionReq{
			UserID:      v.UserID,
			Debt:        v.Debt,
			Status:      v.Status,
			Money:       v.Money,
			CartList:    cartList,
			TotalPrice:  v.TotalPrice,
			Notes:       v.Notes,
			TotalProfit: v.TotalProfit,
		})
	}

	return resp, nil
}
